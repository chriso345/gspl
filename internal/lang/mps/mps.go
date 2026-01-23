package mps

import (
	"bufio"
	"context"
	"io"
	"strconv"
	"strings"

	"github.com/chriso345/gspl/internal/lang"
	"github.com/chriso345/gspl/internal/lang/ast"
	"github.com/chriso345/gspl/lp"
)

// MPSLanguage provides a tiny MPS parser sufficient for the repository examples.
type MPSLanguage struct{}

func (m *MPSLanguage) Name() string { return "mps" }

func parseMPSToLP(r io.Reader, filename string) (*lp.LinearProgram, error) {
	s := bufio.NewScanner(r)
	rows := map[string]rune{}           // name -> type (N,L,G,E)
	cols := map[string]map[string]float64{} // var -> row -> coef
	rhs := map[string]map[string]float64{}  // rhsname -> row -> value
	bounds := []struct{ typ, varname string; val float64 }{}
	var orderVars []string
	desc := filename
	section := ""
	for s.Scan() {
		line := s.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		f := strings.Fields(line)
		if len(f) == 0 {
			continue
		}
		switch strings.ToUpper(f[0]) {
		case "NAME":
			if len(f) > 1 {
				desc = strings.Join(f[1:], " ")
			}
		case "ROWS":
			section = "ROWS"
			continue
		case "COLUMNS":
			section = "COLUMNS"
			continue
		case "RHS":
			section = "RHS"
			continue
		case "BOUNDS":
			section = "BOUNDS"
			continue
		case "ENDATA":
			section = ""
			continue
		}

		switch section {
		case "ROWS":
			// format: type name
			if len(f) >= 2 {
				typ := strings.ToUpper(f[0])
				name := f[1]
				if typ == "N" || typ == "L" || typ == "G" || typ == "E" {
					rows[name] = rune(typ[0])
				}
			}
		case "COLUMNS":
			// format: var  row  val  [row val]
			if len(f) >= 3 {
				varname := f[0]
				if _, ok := cols[varname]; !ok {
					cols[varname] = map[string]float64{}
					orderVars = append(orderVars, varname)
				}
				for i := 1; i+1 < len(f); i += 2 {
					rowname := f[i]
					valStr := f[i+1]
					val, err := strconv.ParseFloat(valStr, 64)
					if err != nil {
						continue
					}
					cols[varname][rowname] = val
				}
			}
		case "RHS":
			// format: name row val [row val]
			if len(f) >= 3 {
				name := f[0]
				if _, ok := rhs[name]; !ok {
					rhs[name] = map[string]float64{}
				}
				for i := 1; i+1 < len(f); i += 2 {
					row := f[i]
					valStr := f[i+1]
					val, err := strconv.ParseFloat(valStr, 64)
					if err != nil {
						continue
					}
					rhs[name][row] = val
				}
			}
		case "BOUNDS":
			// format: typ bname var val
			if len(f) >= 4 {
				typ := f[0]
				varname := f[2]
				valStr := f[3]
				val, err := strconv.ParseFloat(valStr, 64)
				if err == nil {
					bounds = append(bounds, struct{ typ, varname string; val float64 }{typ, varname, val})
				}
			}
		default:
			// ignore
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	// build vars
	vars := []lp.LpVariable{}
	for _, v := range orderVars {
		vars = append(vars, lp.NewVariable(v))
	}
	lprog := lp.NewLinearProgram(desc, vars)

	// objective: find row with type N
	objRow := ""
	for name, t := range rows {
		if t == 'N' {
			objRow = name
			break
		}
	}
	if objRow != "" {
		terms := []lp.LpTerm{}
		for _, v := range vars {
			if coef, ok := cols[v.Name][objRow]; ok {
				terms = append(terms, lp.NewTerm(coef, lp.NewVariable(v.Name)))
			}
		}
		// MPS problems are traditionally minimisation by default
		lprog.AddObjective(lp.LpMinimise, lp.NewExpression(terms))
	}

	// constraints: rows with L,G,E
	for name, t := range rows {
		if t == 'N' {
			continue
		}
		terms := []lp.LpTerm{}
		for _, v := range vars {
			if coef, ok := cols[v.Name][name]; ok {
				terms = append(terms, lp.NewTerm(coef, lp.NewVariable(v.Name)))
			}
		}
		// rhs take first rhs map entry if present
		rhsVal := 0.0
		for _, m := range rhs {
			if val, ok := m[name]; ok {
				rhsVal = val
				break
			}
		}
		var ctype lp.LpConstraintType
		switch t {
		case 'L':
			ctype = lp.LpConstraintLE
		case 'G':
			ctype = lp.LpConstraintGE
		case 'E':
			ctype = lp.LpConstraintEQ
		}
		lprog.AddConstraint(lp.NewExpression(terms), ctype, rhsVal)
	}

	// bounds: implement simple UP and LO parsing from cols of bounds section
	// not implemented in depth; use defaults present in example
	// scan again to pick up bounds (cheap approach)
	if len(bounds) > 0 {
		for _, b := range bounds {
			switch strings.ToUpper(b.typ) {
			case "UP":
				// set upper bound by adding constraint var <= val
				terms := []lp.LpTerm{{Coefficient: 1, Variable: lp.NewVariable(b.varname)}}
				lprog.AddConstraint(lp.NewExpression(terms), lp.LpConstraintLE, b.val)
			case "LO":
				// set lower bound by adding constraint var >= val
				terms := []lp.LpTerm{{Coefficient: 1, Variable: lp.NewVariable(b.varname)}}
				lprog.AddConstraint(lp.NewExpression(terms), lp.LpConstraintGE, b.val)
			}
		}
	}

	return &lprog, nil
}

func (m *MPSLanguage) Parse(ctx context.Context, src io.Reader, opts ...lang.ParseOption) (ast.Node, error) {
	filename := ""
	if f, ok := src.(interface{ Name() string }); ok {
		filename = f.Name()
	}
	lpProg, err := parseMPSToLP(src, filename)
	if err != nil {
		return nil, err
	}
	return &ast.Module{LP: lpProg, Name: lpProg.Description}, nil
}

func New() lang.Language { return &MPSLanguage{} }

func init() { lang.MustRegisterLanguage(New()) }
