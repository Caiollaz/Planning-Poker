package vote

const (
	Confused    = -1 // o card '?'
	NotSelected = -2 // quando o card é desmarcado
	Private     = -3 // este valor retornado quando as cartas ainda não foram reveladas
)

var voteValues = map[int]bool{
	Private:     false, // Privado
	NotSelected: false, // Desmarcar card
	Confused:    false, // Confuso: o card '?'
	0:           true,
	1:           true,
	2:           true,
	3:           true,
	4:           true,
	5:           true,
	8:           true,
	13:          true,
	21:          true,
	34:          true,
	55:          true,
	89:          true,
}
