This is a Go package to generate stroke outlines for cubic Bezier paths.

This package was originally developed for use with Gio, 
but it is not tied to any particular graphics library or GUI toolkit.

Instead of "flattening" cubic curves to sequences of quadratic curves
or line segments like many path-stroking implementations do, 
it uses cubic curves all the way through. 
Linear and quadratic segments are converted to cubic for uniformity.

A path is represented by a `[][]Segment`. 
Each `[]Segment` represents a single contour of the path (open or closed),
with each `Segment` in the contour starting where the previous one ended.
(In the `[][]Segment` returned by `Stroke`, 
starting points and ending points are not guaranteed to coincide exactly.
Callers should fill in the gaps with straight lines, if they occur.)

Graphics libraries generally consider the line's dash pattern to be part of the stroke style,
but this package treats breaking a path into dashes as a separate operation:

```
dashed := stroke.Dash(p, []float32{1, 2}, 0)
outline := stroke.Stroke(dashed, stroke.Options{Width: 2})
```


