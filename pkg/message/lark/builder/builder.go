package builder

import (
	"github.com/larksuite/botframework-go/SDK/message"
	"github.com/larksuite/botframework-go/SDK/protocol"
)

var (
	line1   = 1
	HrThing = &Thing{
		Type: HrThingType,
	}
)

type ThingType int

const (
	TextThingType = iota
	ButtonThingType
	HrThingType
)

type Thing struct {
	Type   ThingType
	Lines  []string
	URLMap map[string]string

	ButtonParamMap    map[string]map[string]string
	ButtonMethodMap   map[string]string
	ButtonCategoryMap map[string]protocol.ButtonStyle
}

type Builder struct {
	Title  string
	Things []*Thing
}

func NewBuilder(title string) *Builder {
	return &Builder{
		Title:  title,
		Things: make([]*Thing, 0),
	}
}

func (bd *Builder) AddThing(things ...*Thing) {
	bd.Things = append(bd.Things, things...)
}

func (bd *Builder) Build() (*protocol.CardForm, error) {
	builder := &message.CardBuilder{}

	builder.SetConfig(
		protocol.ConfigForm{
			MinVersion:     protocol.VersionForm{},
			WideScreenMode: true,
		},
	)

	builder.AddHeader(
		protocol.TextForm{
			Tag:     protocol.PLAIN_TEXT_E,
			Content: &bd.Title,
			Lines:   &line1,
		},
		"",
	)

	builder.AddHRBlock()

	for _, thing := range bd.Things {
		switch thing.Type {
		case TextThingType:
			{
				if len(thing.Lines) == 0 {
					continue
				}

				url := make(map[string]protocol.URLForm)
				for k, v := range thing.URLMap {
					url[k] = protocol.URLForm{
						Url: &v,
					}
				}
				fields := make([]protocol.FieldForm, 0, len(thing.Lines)-1)
				for _, line := range thing.Lines[1:] {
					fields = append(fields, protocol.FieldForm{
						Short: false,
						Text:  *message.NewMDText(line, nil, nil, url),
					})
				}

				builder.AddDIVBlock(
					message.NewMDText(thing.Lines[0], nil, nil, url),
					fields,
					nil,
				)
			}
		case HrThingType:
			{
				builder.AddHRBlock()
			}
		case ButtonThingType:
			{
				elems := make([]protocol.ActionElement, 0, len(thing.Lines))
				for _, line := range thing.Lines {
					elems = append(
						elems,
						message.NewButton(
							message.NewMDText(line, nil, nil, nil),
							nil,
							nil,
							thing.ButtonParamMap[line],
							thing.ButtonCategoryMap[line],
							nil,
							thing.ButtonMethodMap[line],
						),
					)
				}

				builder.AddActionBlock(elems)
			}
		}
	}

	if len(bd.Things) != 0 && bd.Things[len(bd.Things)-1].Type != HrThingType {
		builder.AddHRBlock()
	}

	// noteTextContent := "有下标看起来就很专业®√"
	// builder.AddNoteBlock(
	// 	[]protocol.BaseElement{
	// 		&protocol.TextForm{
	// 			Tag:     protocol.PLAIN_TEXT_E,
	// 			Content: &noteTextContent,
	// 			Lines:   nil,
	// 		},
	// 	})

	card, err := builder.BuildForm()
	if err != nil {
		return nil, err
	}

	return card, nil
}
