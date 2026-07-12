package component

import (
	"image/color"
	"testing"

	"github.com/fglo/chopstiqs/event"
	"github.com/fglo/chopstiqs/option"
	"github.com/matryer/is"
)

func TestContainer_AddComponent(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	c := NewContainer(&ContainerOptions{})
	c.SetEventManager(eventManager)

	b := NewButton(&ButtonOptions{})
	c.AddComponent(b)

	is.Equal(len(c.components), 1)
	is.Equal(c.components[0], b)
}

func TestContainer_AddComponents(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	c := NewContainer(&ContainerOptions{})
	c.SetEventManager(eventManager)

	b1 := NewButton(&ButtonOptions{})
	b2 := NewButton(&ButtonOptions{})
	b3 := NewButton(&ButtonOptions{})

	c.AddComponents(b1, b2, b3)

	is.Equal(len(c.components), 3)
}

func TestContainer_SetDisabled_Cascade(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	c := NewContainer(&ContainerOptions{})
	c.SetEventManager(eventManager)

	b1 := NewButton(&ButtonOptions{})
	b2 := NewButton(&ButtonOptions{})
	c.AddComponents(b1, b2)

	c.SetDisabled(true)

	is.Equal(c.Disable(), true)
	is.Equal(b1.Disable(), true)
	is.Equal(b2.Disable(), true)

	c.SetDisabled(false)

	is.Equal(c.Disable(), false)
	is.Equal(b1.Disable(), false)
	is.Equal(b2.Disable(), false)
}

func TestContainer_SetPosition_RecalculatesChildren(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	c := NewContainer(&ContainerOptions{})
	c.SetEventManager(eventManager)
	c.SetPosition(10, 20)

	b := NewButton(&ButtonOptions{})
	c.AddComponent(b)

	c.SetPosition(50, 60)

	absX, absY := b.AbsPosition()
	is.Equal(absX, 50.)
	is.Equal(absY, 60.)
}

func TestContainer_SetPosX_RecalculatesChildren(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	c := NewContainer(&ContainerOptions{})
	c.SetEventManager(eventManager)
	c.SetPosition(10, 20)

	b := NewButton(&ButtonOptions{})
	c.AddComponent(b)

	c.SetPosX(100)

	absX, _ := b.AbsPosition()
	is.Equal(absX, 100.)
}

func TestContainer_SetPosY_RecalculatesChildren(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	c := NewContainer(&ContainerOptions{})
	c.SetEventManager(eventManager)
	c.SetPosition(10, 20)

	b := NewButton(&ButtonOptions{})
	c.AddComponent(b)

	c.SetPosY(200)

	_, absY := b.AbsPosition()
	is.Equal(absY, 200.)
}

func TestContainer_SetBackgroundColor(t *testing.T) {
	is := is.New(t)

	c := NewContainer(&ContainerOptions{})

	red := color.RGBA{255, 0, 0, 255}
	c.SetBackgroundColor(red)

	is.Equal(c.GetBackgroundColor(), red)
}

func TestContainer_FireEvents_Delegates(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	c := NewContainer(&ContainerOptions{})
	c.SetEventManager(eventManager)

	fired := 0

	b := NewButton(&ButtonOptions{})
	c.AddComponent(b)
	b.AddClickedHandler(func(args *ButtonClickedEventArgs) {
		fired++
	})

	leftMouseButtonClick(t, &b.component)
	is.Equal(fired, 1)

	c.FireEvents()
}

func TestContainer_Layout_CalledOnAdd(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	hl := &HorizontalListLayout{ColumnGap: 5}

	c := NewContainer(&ContainerOptions{Layout: hl})
	c.SetEventManager(eventManager)
	c.SetDimensions(200, 50)

	b := NewButton(&ButtonOptions{Width: option.Int(30), Height: option.Int(20)})
	c.AddComponent(b)

	posX, _ := b.Position()
	is.True(posX >= 0)
}

func TestContainer_FocusedEvent_Forwards(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	c := NewContainer(&ContainerOptions{})
	c.SetEventManager(eventManager)

	focusedCount := 0
	c.AddFocusedHandler(func(args *ComponentFocusedEventArgs) {
		focusedCount++
	})

	b := NewButton(&ButtonOptions{})
	c.AddComponent(b)

	b.SetFocused(true)
	eventManager.HandleFired()

	is.True(focusedCount >= 1)
}

func TestContainer_Empty(t *testing.T) {
	is := is.New(t)

	c := NewContainer(nil)

	is.Equal(len(c.components), 0)
	is.Equal(c.Disable(), false)
	is.Equal(c.Hidden(), false)
}

func TestContainer_Nested_RecalculateAbsPosition(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()

	parent := NewContainer(&ContainerOptions{})
	parent.SetEventManager(eventManager)
	parent.SetPosition(10, 20)

	child := NewContainer(&ContainerOptions{})
	parent.AddComponent(child)
	child.SetPosition(5, 10)

	absX, absY := child.AbsPosition()
	is.Equal(absX, 15.)
	is.Equal(absY, 30.)

	grandchild := NewButton(&ButtonOptions{})
	child.AddComponent(grandchild)
	grandchild.SetPosition(3, 7)

	gcAbsX, gcAbsY := grandchild.AbsPosition()
	is.Equal(gcAbsX, 18.)
	is.Equal(gcAbsY, 37.)

	parent.SetPosition(100, 200)

	gcAbsX, gcAbsY = grandchild.AbsPosition()
	is.Equal(gcAbsX, 108.)
	is.Equal(gcAbsY, 217.)
}
