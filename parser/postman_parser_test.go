package test_report

import "testing"

func TestFindItemId(t *testing.T) {
	scenario := []PostmanItem{
		{
			Name: "Folder",
			Id:   "1",
			Item: []PostmanItem{
				{
					Name: "Request",
					Id:   "2",
					Item: nil,
				},
			},
		},
		{
			Name: "Folder2",
			Id:   "3",
			Item: []PostmanItem{
				{
					Name: "Subfolder",
					Id:   "4",
					Item: []PostmanItem{
						{
							Name: "Request",
							Id:   "5",
							Item: nil,
						},
					},
				},
			},
		},
	}
	result := findPathToId("2", "", scenario)
	if result != "Folder" {
		t.Fatal(result)
	}
	result = findPathToId("5", "", scenario)
	if result != "Folder2/Subfolder" {
		t.Fatal(result)
	}
}
