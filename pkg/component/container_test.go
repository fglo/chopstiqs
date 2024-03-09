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
			xmlContents: "<container layout=\"name: verticalList, rowGap: 5\" width=\"129\" height=\"170\" padding=\"5,5,5,5\"></container>",
			wantErr:     false,
			want: &Container{
				layout: &VerticalListLayout{RowGap: 5},
				component: component{
					width:   129,
					height:  170,
					padding: Padding{5, 5, 5, 5},
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
