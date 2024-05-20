package services

type TreatmentFunc func(lastResult, currentResult interface{}, attributeId string, filters []string) ([]interface{}, error)

func TreatmentFactory(differential string) TreatmentFunc {
	switch differential {
	case "new":
		return NewTreatment
	case "delete":
		return DeleteTreatment
	case "update":
		return UpdateTreatment
	default:
		return nil
	}
}
