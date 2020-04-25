package geodata

type UseCase interface {
	CheckGeoPos(address string) (bool, error)
}
