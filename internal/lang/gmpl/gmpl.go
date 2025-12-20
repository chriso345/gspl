package gmpl

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/chriso345/gspl/internal/lang"
	"github.com/chriso345/gspl/internal/lang/ast"
	"github.com/chriso345/gspl/lp"
)

// GMPLLanguage is a small GMPL adapter that supports very basic subsets used in examples.
type GMPLLanguage struct{}

func (g *GMPLLanguage) Name() string { return "gmpl" }

// parseExpression parses simple expressions like "3*x1 + 2*x2" or "x1 + x2" and returns pairs of (coef, varname)
func parseExpression(s string) ([]struct {
	coef float64
	name string
}, error) {
	s = strings.TrimSpace(s)
	// replace '-' with '+-' to simplify splitting
	s2 := strings.ReplaceAll(s, "-", "+-")
	tokens := strings.Split(s2, "+")
	var out []struct {
		coef float64
		name string
	}
	for _, tok := range tokens {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		}
		coef := 1.0
		name := ""
		if strings.Contains(tok, "*") {
			parts := strings.Split(tok, "*")
			cstr := strings.TrimSpace(parts[0])
			nstr := strings.TrimSpace(parts[1])
			switch cstr {
			case "", "+":
				coef = 1.0
			case "-":
				coef = -1.0
			default:
				c, err := strconv.ParseFloat(cstr, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid coefficient %q: %w", cstr, err)
				}
				coef = c
			}
			name = strings.TrimSuffix(nstr, ";")
		} else {
			p := strings.Fields(tok)
			if len(p) == 1 {
				if strings.HasPrefix(p[0], "-") {
					coef = -1.0
					name = strings.TrimPrefix(p[0], "-")
				} else {
					name = strings.TrimSuffix(p[0], ";")
				}
			} else if len(p) == 2 {
				c, err := strconv.ParseFloat(p[0], 64)
				if err != nil {
					return nil, fmt.Errorf("invalid coefficient %q: %w", p[0], err)
				}
				coef = c
				name = strings.TrimSuffix(p[1], ";")
			} else {
				name = strings.TrimSuffix(p[len(p)-1], ";")
			}
		}
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		out = append(out, struct {
			coef float64
			name string
		}{coef, name})
	}
	return out, nil
}

func parseReaderToLP(r io.Reader, filename string) (*lp.LinearProgram, error) {
	s := bufio.NewScanner(r)
	vars := []string{}
	objectiveSense := lp.LpMinimise
	objectiveExpr := ""
	constraints := []struct {
		lhs string
		typ string
		rhs float64
	}{}

	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// strip trailing ;
		line = strings.TrimSuffix(line, ";")
		lower := strings.ToLower(line)
		if strings.HasPrefix(lower, "var ") {
			// var x1 >= 0;  -> extract variable name
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				vars = append(vars, parts[1])
			}
			continue
		}
		if strings.HasPrefix(lower, "maximize") || strings.HasPrefix(lower, "minimize") {
			if strings.HasPrefix(lower, "maximize") {
				objectiveSense = lp.LpMaximise
			}
			if idx := strings.Index(line, ":"); idx >= 0 {
				objectiveExpr = strings.TrimSpace(line[idx+1:])
			} else if idx := strings.Index(line, " "); idx >= 0 {
				objectiveExpr = strings.TrimSpace(line[idx+1:])
			}
			continue
		}
		if strings.HasPrefix(lower, "subject to") {
			rest := strings.TrimSpace(line[len("subject to"):])
			if idx := strings.Index(rest, ":"); idx >= 0 {
				rest = strings.TrimSpace(rest[idx+1:])
			}
			if strings.Contains(rest, "<=") {
				parts := strings.SplitN(rest, "<=", 2)
				lhs := strings.TrimSpace(parts[0])
				rhsStr := strings.TrimSpace(parts[1])
				rhs, err := strconv.ParseFloat(rhsStr, 64)
				if err != nil {
					return nil, err
				}
				constraints = append(constraints, struct {
					lhs string
					typ string
					rhs float64
				}{lhs, "<=", rhs})
				continue
			}
			if strings.Contains(rest, ">=") {
				parts := strings.SplitN(rest, ">=", 2)
				lhs := strings.TrimSpace(parts[0])
				rhsStr := strings.TrimSpace(parts[1])
				rhs, err := strconv.ParseFloat(rhsStr, 64)
				if err != nil {
					return nil, err
				}
				constraints = append(constraints, struct {
					lhs string
					typ string
					rhs float64
				}{lhs, ">=", rhs})
				continue
			}
			if strings.Contains(rest, "=") {
				parts := strings.SplitN(rest, "=", 2)
				lhs := strings.TrimSpace(parts[0])
				rhsStr := strings.TrimSpace(parts[1])
				rhs, err := strconv.ParseFloat(rhsStr, 64)
				if err != nil {
					return nil, err
				}
				constraints = append(constraints, struct {
					lhs string
					typ string
					rhs float64
				}{lhs, "=", rhs})
				continue
			}
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	// build LP
	lpVars := make([]lp.LpVariable, len(vars))
	for i, v := range vars {
		lpVars[i] = lp.NewVariable(v)
	}
	lprog := lp.NewLinearProgram(filename, lpVars)
	// objective
	if objectiveExpr != "" {
		parsed, err := parseExpression(objectiveExpr)
		if err != nil {
			return nil, err
		}
		terms := []lp.LpTerm{}
		for _, p := range parsed {
			varObj := lp.NewVariable(p.name)
			terms = append(terms, lp.NewTerm(p.coef, varObj))
		}
		lprog.AddObjective(objectiveSense, lp.NewExpression(terms))
	}
	// constraints
	for _, c := range constraints {
		parsed, err := parseExpression(c.lhs)
		if err != nil {
			return nil, err
		}
		terms := []lp.LpTerm{}
		for _, p := range parsed {
			varObj := lp.NewVariable(p.name)
			terms = append(terms, lp.NewTerm(p.coef, varObj))
		}
		var conType lp.LpConstraintType
		switch c.typ {
		case "<=":
			conType = lp.LpConstraintLE
		case ">=":
			conType = lp.LpConstraintGE
		case "=":
			conType = lp.LpConstraintEQ
		}
		lprog.AddConstraint(lp.NewExpression(terms), conType, c.rhs)
	}

	return &lprog, nil
}

func (g *GMPLLanguage) Parse(ctx context.Context, src io.Reader, opts ...lang.ParseOption) (ast.Node, error) {
	// Attempt to detect filename for description
	filename := ""
	if f, ok := src.(interface{ Name() string }); ok {
		filename = f.Name()
	}
	lpProg, err := parseReaderToLP(src, filename)
	if err != nil {
		return nil, err
	}
	m := &ast.Module{LP: lpProg, Name: lpProg.Description}
	return m, nil
}

func New() lang.Language { return &GMPLLanguage{} }

func init() { lang.MustRegisterLanguage(New()) }
