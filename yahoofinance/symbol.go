package yahoofinance

var (
	SymbolXEQT = newSymbol("XEQT.TO")
)

type Symbol struct {
	id string
}

func newSymbol(id string) Symbol {
	return Symbol{id: id}
}

func (s Symbol) String() string {
	return s.id
}
