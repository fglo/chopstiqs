package component

import (
	"encoding/xml"
	"testing"
)

func Test_UnmarshalXML(t *testing.T) {
	tests := []struct {
		name        string
		xmlContents string
		wantErr     bool
		want        *Container
	}{
		{
			name:        "Unmarshal VerticalListLayout",
			xmlContents: "<container layout=\"name: verticalList, rowGap: 5\" width=\"129\" height=\"170\" padding=\"1,2,3,4\"></container>",
			wantErr:     false,
			want: &Container{
				layout: &VerticalListLayout{RowGap: 5},
				component: component{
					width:   129,
					height:  170,
					padding: Padding{Top: 1, Right: 2, Bottom: 3, Left: 4},
				},
			},
		},
		{
			name:        "Unmarshal HorizontalListLayout",
			xmlContents: "<container layout=\"name: horizontalList, columnGap: 5\" width=\"129\" height=\"170\" padding=\"1,2,3,4\"></container>",
			wantErr:     false,
			want: &Container{
				layout: &HorizontalListLayout{ColumnGap: 5},
				component: component{
					width:   129,
					height:  170,
					padding: Padding{Top: 1, Right: 2, Bottom: 3, Left: 4},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Container{}

			err := xml.Unmarshal([]byte(tt.xmlContents), c)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Container.UnmarshalXML() error = %v, wantErr %v", err, tt.wantErr)
			}

			CompareContainers(t, c, tt.want)
		})
	}
}
