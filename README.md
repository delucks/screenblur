# screenblur

This utility takes a screenshot and applies a blur to the screenshot, quickly. That's it. This special-purpose tool is designed for use with window managers (like `i3`) that package a screen locking utility, to reduce the time between pressing your lock shortcut and the screen locking. The blur algorithm used is [Stackblur](https://github.com/esimov/stackblur-go) which is faster than a Gaussian blur.

# Why?

It's almost become convention in some unix circles to lock your desktop with a blurred screenshot of the desktop. Many distributions that use tiling window managers by default include a utility which takes a screenshot, applies a blur, then locks the desktop using that image. Unfortunately, those utilities are generally very slow. Here's one, included by default with Manjaro's i3 flavor.

```
$ time blurlock
real    0m13.680s
user    0m23.629s
sys     0m0.380s
```

The implementation uses ImageMagick, which isn't particularly speedy when invoked from the CLI tools:

```
#!/bin/bash
# /usr/bin/blurlock

# take screenshot
import -window root /tmp/screenshot.png

# blur it
convert /tmp/screenshot.png -blur 0x5 /tmp/screenshotblur.png
rm /tmp/screenshot.png

# lock the screen
i3lock -i /tmp/screenshotblur.png

# sleep 1 adds a small delay to prevent possible race conditions with suspend
sleep 1

exit 0
```

(All benchmarks comment out the "i3lock" and "sleep" calls to get a more accurate measurement of the image processing portion)

Using a different image capture program like `feh` or `maim` speeds things up drastically, but we're still looking at the ~5s range on average. `screenblur` is an effort to improve the state of the art by using special-purpose code written in a compiled language.

Here's a benchmark comparing these three options. `./block` uses `maim` instead of `import` and `./blurlock` is the above implementation.

```
$ hyperfine './block' './blurlock' './screenblur'
Benchmark #1: ./block
  Time (mean ± σ):      6.758 s ±  0.759 s    [User: 17.217 s, System: 0.260 s]
  Range (min … max):    6.000 s …  7.982 s    10 runs

Benchmark #2: ./blurlock
  Time (mean ± σ):     34.917 s ± 10.773 s    [User: 45.620 s, System: 0.388 s]
  Range (min … max):   14.139 s … 41.399 s    10 runs

Benchmark #3: ./screenblur
  Time (mean ± σ):      1.559 s ±  0.133 s    [User: 1.400 s, System: 0.134 s]
  Range (min … max):    1.351 s …  1.707 s    10 runs

Summary
  './screenblur' ran
    4.34 ± 0.61 times faster than './block'
   22.40 ± 7.17 times faster than './blurlock'
```
