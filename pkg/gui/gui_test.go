package gui

import (
	"encoding/xml"
	"testing"

	"github.com/fglo/chopstiqs/pkg/component"
	"github.com/fglo/chopstiqs/pkg/option"
)

func Test_UnmarshalXML(t *testing.T) {
	tests := []struct {
		name        string
		xmlContents string
		wantErr     bool
		want        GUI
	}{
		{
			name:        "Unmarshal VerticalListLayout",
			xmlContents: "<container layout=\"name: verticalList, rowGap: 5\" width=\"129\" height=\"170\" padding=\"5,5,5,5\"></container>",
			wantErr:     false,
			want: GUI{
				rootContainer: component.NewContainer(&component.ContainerOptions{
					Layout:  &component.VerticalListLayout{RowGap: 5},
					Padding: component.NewPadding(5, 5, 5, 5),
					Width:   option.Int(129),
					Height:  option.Int(170),
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gui := &GUI{}

			err := xml.Unmarshal([]byte(tt.xmlContents), gui)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Container.UnmarshalXML() error = %v, wantErr %v", err, tt.wantErr)
			}

			component.CompareContainers(t, gui.rootContainer, tt.want.rootContainer)
		})
	}
}
