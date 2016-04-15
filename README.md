# RiGO
Implementation of the RenderMan Interface for the Go programming language. This is currently 
based on Pixar's RenderMan Specification version 3.2.1 (November 2005). This implementation 
is still under *active development*, so *expect* holes and bugs. 

[Online Documentation](https://godoc.org/github.com/mae-global/rigo)

Install with:

    go get github.com/mae-global/rigo

Quick example usage; outputting a Unit Cube to a RIB Entity file. 

![Pipeline](documentation/pipe.png)

```go
/* create a function to record the duration between RiBegin and RiEnd calls */
type MyTimer struct {
	start time.Time
	finish time.Time
}

func (t MyTimer) Name() string {
	return "mytimer"
}

func (t *MyTimer) Took() time.Duration {
	return t.finish.Sub(t.start)
}

func (t *MyTimer) Write(name RtName,list []Rter,info Info) *Result {
	switch string(name) {
		case "Begin","RiBegin":
			t.start = time.Now()
			t.finish = t.start
		break
		case "End","RiEnd":
			t.finish = time.Now()
		break
	}
	return Done()
}

/* Construct a pipeline, including our timer, piping RIB output to file */
pipe := NewPipe()
pipe.Append(&MyTimer{}).Append(&PipeToFile{})

ctx := NewEntity(pipe)

/* Do all our Ri calls */
ctx.Begin("unitcube.rib")
ctx.AttributeBegin("begin unit cube")
ctx.Attribute("identifier", RtToken("name"), RtToken("unitcube"))
ctx.Bound(RtBound{-.5, .5, -.5, .5, -.5, .5})
ctx.TransformBegin()

ctx.Comment("far face")
ctx.Polygon(4, RtToken("P"), RtFloatArray{.5, .5, .5, -.5, .5, .5, -.5, -.5, .5, .5, -.5, .5})
ctx.Rotate(90, 0, 1, 0)

ctx.Comment("right face")
ctx.Polygon(4, RtToken("P"), RtFloatArray{.5, .5, .5, -.5, .5, .5, -.5, -.5, .5, .5, -.5, .5})
ctx.Rotate(90, 0, 1, 0)

ctx.Comment("near face")
ctx.Polygon(4, RtToken("P"), RtFloatArray{.5, .5, .5, -.5, .5, .5, -.5, -.5, .5, .5, -.5, .5})
ctx.Rotate(90, 0, 1, 0)

ctx.Comment("left face")
ctx.Polygon(4, RtToken("P"), RtFloatArray{.5, .5, .5, -.5, .5, .5, -.5, -.5, .5, .5, -.5, .5})

ctx.TransformEnd()
ctx.TransformBegin()

ctx.Comment("bottom face")
ctx.Rotate(90, 1, 0, 0)
ctx.Polygon(4, RtToken("P"), RtFloatArray{.5, .5, .5, -.5, .5, .5, -.5, -.5, .5, .5, -.5, .5})

ctx.TransformEnd()
ctx.TransformBegin()

ctx.Comment("top face")
ctx.Rotate(-90, 1, 0, 0)
ctx.Polygon(4, RtToken("P"), RtFloatArray{.5, .5, .5, -.5, .5, .5, -.5, -.5, .5, .5, -.5, .5})

ctx.TransformEnd()
ctx.AttributeEnd("end unit cube")
ctx.End()	
		
/* grab our timer back and print the duration */
p = pipe.GetByName(MyTimer{}.Name())
t,_ := p.(*MyTimer)
	
fmt.Printf("took %s\n",t.Took())
```	

RIB output of *unitcube.rib* is thus :-

```
##RenderMan RIB-Structure 1.1 Entity
AttributeBegin #begin unit cube
	Attribute "identifier" "name" "unitcube"
	Bound [-.5 .5 -.5 .5 -.5 .5]
	TransformBegin 
		# far face
		Polygon "P" [.5 .5 .5 -.5 .5 .5 -.5 -.5 .5 .5 -.5 .5]
		Rotate 90. 0 1. 0
		# right face
		Polygon "P" [.5 .5 .5 -.5 .5 .5 -.5 -.5 .5 .5 -.5 .5]
		Rotate 90. 0 1. 0
		# near face
		Polygon "P" [.5 .5 .5 -.5 .5 .5 -.5 -.5 .5 .5 -.5 .5]
		Rotate 90. 0 1. 0
		# left face
		Polygon "P" [.5 .5 .5 -.5 .5 .5 -.5 -.5 .5 .5 -.5 .5]
	TransformEnd 
	TransformBegin 
		# bottom face
		Rotate 90. 1. 0 0
		Polygon "P" [.5 .5 .5 -.5 .5 .5 -.5 -.5 .5 .5 -.5 .5]
	TransformEnd 
	TransformBegin 
		# top face
		Rotate -90. 1. 0 0
		Polygon "P" [.5 .5 .5 -.5 .5 .5 -.5 -.5 .5 .5 -.5 .5]
	TransformEnd 
AttributeEnd #end unit cube
```
We can remove duplicate work, by using the fragments interface.

```go
pipe := DefaultFilePipe()

ctx := New(pipe,nil)
ctx.Begin("output/orangeball.rib")
ctx.ArchiveRecord("structure","Scene Orange Ball")
ctx.ArchiveRecord("structure","Frames 5")
		
frag := NewFragment("orangeball_fragment")
/* grab the Renderman Interface from the fragment */
fri := frag.Ri()

fri.Format(640,480,-1)
fri.ShadingRate(1)
fri.Projection(Perspective,RtString("fov"),RtInt(30))
fri.FrameAspectRatio(1.33)
fri.Identity()
fri.LightSource("distantlight",RtInt(1))
fri.Translate(0,0,5)
fri.WorldBegin()
fri.Identity()
fri.AttributeBegin()
fri.Color(RtColor{1.0,0.6,0.0})
fri.Surface("plastic",RtString("Ka"),RtFloat(1),RtString("Kd"),RtFloat(0.5),
										  RtString("Ks"),RtFloat(1),RtString("roughness"),RtFloat(0.1))
fri.TransformBegin()
fri.Rotate(90,1,0,0)
fri.Sphere(1,-1,1,360)
fri.TransformEnd()
fri.AttributeEnd()
fri.WorldEnd()


for frame := 1; frame <= 5; frame++ {
	ctx.FrameBegin(RtInt(frame))
	ctx.Display(RtToken(fmt.Sprintf("orange_%03d.tif",frame)),"file","rgba")		

	frag.Replay(ctx)

	ctx.FrameEnd()
}

ctx.End()
```
```
##RenderMan RIB-Structure 1.1
##Scene Orange Ball
##Creator rigo-0
##CreationDate 2016-04-15 13:06:05.461257304 +0100 BST
##For mae
##Frames 5
FrameBegin 1
	Display "orange_001.tif" "file" "rgba"
	Format 640 480 -1
	ShadingRate 1
	Projection "perspective" "fov" 30
	FrameAspectRatio 1.33
	Identity 
	LightSource "distantlight" "0"
	Translate 0 0 5
	WorldBegin 
		Identity 
		AttributeBegin 
			Color [1 .6 0]
			Surface "plastic" "Ka" 1 "Kd" .5 "Ks" 1 "roughness" .1
			TransformBegin 
				Rotate 90 1 0 0
				Sphere 1 -1 1 360
			TransformEnd 
		AttributeEnd 
	WorldEnd 
FrameEnd 
FrameBegin 2
	Display "orange_002.tif" "file" "rgba"
	Format 640 480 -1
	ShadingRate 1
	Projection "perspective" "fov" 30
	FrameAspectRatio 1.33
	Identity 
	LightSource "distantlight" "0"
	Translate 0 0 5
	WorldBegin 
		Identity 
		AttributeBegin 
			Color [1 .6 0]
			Surface "plastic" "Ka" 1 "Kd" .5 "Ks" 1 "roughness" .1
			TransformBegin 
				Rotate 90 1 0 0
				Sphere 1 -1 1 360
			TransformEnd 
		AttributeEnd 
	WorldEnd 
FrameEnd 
FrameBegin 3
	Display "orange_003.tif" "file" "rgba"
	Format 640 480 -1
	ShadingRate 1
	Projection "perspective" "fov" 30
	FrameAspectRatio 1.33
	Identity 
	LightSource "distantlight" "0"
	Translate 0 0 5
	WorldBegin 
		Identity 
		AttributeBegin 
			Color [1 .6 0]
			Surface "plastic" "Ka" 1 "Kd" .5 "Ks" 1 "roughness" .1
			TransformBegin 
				Rotate 90 1 0 0
				Sphere 1 -1 1 360
			TransformEnd 
		AttributeEnd 
	WorldEnd 
FrameEnd 
FrameBegin 4
	Display "orange_004.tif" "file" "rgba"
	Format 640 480 -1
	ShadingRate 1
	Projection "perspective" "fov" 30
	FrameAspectRatio 1.33
	Identity 
	LightSource "distantlight" "0"
	Translate 0 0 5
	WorldBegin 
		Identity 
		AttributeBegin 
			Color [1 .6 0]
			Surface "plastic" "Ka" 1 "Kd" .5 "Ks" 1 "roughness" .1
			TransformBegin 
				Rotate 90 1 0 0
				Sphere 1 -1 1 360
			TransformEnd 
		AttributeEnd 
	WorldEnd 
FrameEnd 
FrameBegin 5
	Display "orange_005.tif" "file" "rgba"
	Format 640 480 -1
	ShadingRate 1
	Projection "perspective" "fov" 30
	FrameAspectRatio 1.33
	Identity 
	LightSource "distantlight" "0"
	Translate 0 0 5
	WorldBegin 
		Identity 
		AttributeBegin 
			Color [1 .6 0]
			Surface "plastic" "Ka" 1 "Kd" .5 "Ks" 1 "roughness" .1
			TransformBegin 
				Rotate 90 1 0 0
				Sphere 1 -1 1 360
			TransformEnd 
		AttributeEnd 
	WorldEnd 
FrameEnd 
```
An example light handler generator, which generates unique names so that lights can be tracked more easily. 

```go
pipe := DefaultFilePipe()
	
/* use a custom unique generator with a prefix for the light handles */
lights := NewPrefixLightUniqueGenerator("light_")
	
ctx := NewCustom(pipe,lights,nil,&Configuration{PrettyPrint:true})
ctx.Begin("output/simple.rib")
ctx.Display("sphere.tif","file","rgb")
ctx.Format(320,240,1)
ctx.Projection(Perspective,RtString("fov"),RtFloat(30))
ctx.Translate(0,0,6)
ctx.WorldBegin()
ctx.LightSource("ambientlight",RtString("intensity"),RtFloat(0.5))
ctx.LightSource("distantlight",RtString("intensity"),RtFloat(1.2),RtString("form"),RtIntArray{0,0,-6},RtString("to"),RtIntArray{0,0,0})
ctx.Color(RtColor{1,0,0})
ctx.Sphere(1,-1,1,360)
ctx.WorldEnd()
ctx.End()
```

```
##RenderMan RIB-Structure 1.1
Display "sphere.tif" "file" "rgb"
Format 320 240 1
Projection "perspective" "fov" 30
Translate 0 0 6
WorldBegin 
	LightSource "ambientlight" "light_09c84b71" "intensity" .2
	LightSource "distantlight" "light_64f4dfbf" "intensity" 1.2 "form" [0 0 -6] "to" [0 0 0]
	Color [1 0 0]
	Sphere 1 -1 1 360
WorldEnd 
```





##Roadmap

- [x] Basic RIB pipe
- [ ] Complete RenderMan Interface
- [ ] Stdout/buffer wrapper around io.Writer interface
- [ ] Complete Error checking for each Ri Call
  - [x] Basic Error checking
	- [ ] Sanity checking
	- [ ] Per call checking
	- [ ] Parameterlist checking
- [ ] RIB parser
- [ ] Call wrapping for Ri[call]Begin/Ri[call]End pairs
- [x] Call Fragments 
- [ ] Documentation/Examples


###Information

RenderMan Interface Specification is Copyright © 2005 Pixar.
RenderMan © is a registered trademark of Pixar.

