package ri

import (
	"fmt"
)

/* MakeTexture Convert an image in a standard picture file whise name is picturename */
func (ctx *Context) MakeTexture(picturename, texturename RtString, swrap, twrap RtToken, filterfunc RtFilterFunc, swidth, twidth RtFloat, parameterlist ...Rter) error {

	var out = []Rter{picturename, texturename, swrap, twrap, filterfunc, swidth, twidth}
	out = append(out, parameterlist...)

	return ctx.writef("MakeTexture", out...)
}

/* MakeLatLongEnvironment */
func (ctx *Context) MakeLatLongEnvironment(picturename, texturename RtString, filterfunc RtFilterFunc, swidth, twidth RtFloat, parameterlist ...Rter) error {

	var out = []Rter{picturename, texturename, filterfunc, swidth, twidth}
	out = append(out, parameterlist...)

	return ctx.writef("MakeLatLongEnvironment", out...)
}

/* MakeCubeFaceEnviroment */
func (ctx *Context) MakeCubeFaceEnvironment(px, nx, py, ny, pz, nz, texturename RtString, fov RtFloat, filterfunc RtFilterFunc, swidth, twidth RtFloat, parameterlist ...Rter) error {

	var out = []Rter{px, nx, py, ny, pz, nz, texturename, fov, filterfunc, swidth, twidth}
	out = append(out, parameterlist...)

	return ctx.writef("MakeCubeFaceEnvironment", out...)
}

/* MakeShadow */
func (ctx *Context) MakeShadow(picturename, texturename RtString, parameterlist ...Rter) error {

	var out = []Rter{picturename, texturename}
	out = append(out, parameterlist...)

	return ctx.writef("MakeShadow", out...)
}

/* ArchiveRecord */
func (ctx *Context) ArchiveRecord(typeof RtToken, format RtString, args ...interface{}) error {

	var err error

	switch string(typeof) {
	case "comment":

		err = ctx.writef("#", RtName(fmt.Sprintf(string(format), args...)))
		break
	case "structure":

		err = ctx.writef("##", RtName(fmt.Sprintf(string(format), args...)))
		break
	case "verbatim":

		err = ctx.writef("Verbatim", RtName(fmt.Sprintf(string(format), args...)))
		break
	}

	return err
}

/* ReadArchive */
func (ctx *Context) ReadArchive(name RtToken, callback RtArchiveCallbackFunc, parameterlist ...Rter) error {

	var out = []Rter{name}
	out = append(out, parameterlist...)

	return ctx.writef("ReadArchive", out...)
}