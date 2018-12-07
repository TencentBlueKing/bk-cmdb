package metadata

import (
	"encoding/json"
	"testing"

	"configcenter/src/common/condition"
)

func TestSeachInputConvSearchConds(t *testing.T) {

	s := &SearchInput{
		Fields: []string{"bk_host_id"},
		SortArr: []SearchSort{
			SearchSort{
				Field: "bk_host_id",
				IsDsc: true,
			},
		},
		Limit: &SearchLimit{
			Offset: 0,
			Limit:  100,
		},
		Condition: []SearchInputConditionItem{

			SearchInputConditionItem{
				Fields:   "bk_host_id",
				Operator: condition.BKDBGT,
				Value:    0,
			},
			SearchInputConditionItem{
				Fields:   "$OR",
				Operator: condition.BKDBGT,
				Value: SearchInputConditionItem{
					Fields:   "bk_host_id",
					Operator: condition.BKDBGT,
					Value:    0,
				},
			},
		},
	}

	searchConds := s.ToSearchCondition()
	sBytes, _ := json.Marshal(searchConds)
	result := `{"fields":["bk_host_id"],"condition":{"$OR":{"$gt":{"bk_host_id":{"$gt":0}}},"bk_host_id":{"$gt":0}}}`

	if result != string(sBytes) {
		t.Errorf("result equal")
		return
	}
}

func TestSearchInputJSONConvSearchConds(t *testing.T) {

	str := `{"fields":["bk_host_id"],"limit":{"start":0,"limit":100},"sort":[{"is_dsc":true,"field":"bk_host_id"}],"condition":[{"field":"bk_host_id","operator":"$gt","value":0},{"field":"$OR","operator":"$gt","value":{"field":"bk_host_id","operator":"$gt","value":0}}]}`
	s := &SearchInput{}
	err := json.Unmarshal([]byte(str), s)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	searchConds := s.ToSearchCondition()
	sBytes, _ := json.Marshal(searchConds)
	result := `{"fields":["bk_host_id"],"condition":{"$OR":{"$gt":{"bk_host_id":{"$gt":0}}},"bk_host_id":{"$gt":0}}}`

	if result != string(sBytes) {
		t.Errorf("result equal")
		return
	}

}
