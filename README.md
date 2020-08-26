Tools for i3wm
==

cmd/scratches
--
This tool allows you to display your scratchpad windows in polybar. When you click their names in bar, you toggle them.
When scratchpad window is located on currently visible workspace its name is highlighted.

*scratches* truncates window title up to 20 characters by default. If you need to disable it just pass `-trunacte 0` argument.
If you have music player in scratchpad you can turn truncate off and see full now playing title, for example.

![](./img/scratches.png)

For each of your scratchpads add this section in polybar config:
```
[module/scratch-SCRATCHPAD_NAME]
type = custom/script
exec = scratches -s SCRATCHPAD_NAME
label = %output:0:20:...%
format =  <label>
format-foreground = ${colors.foreground-alt}
click-left = scratches -s SCRATCHPAD_NAME -op toggle
```

**Command line arguments**

```
Usage of scratches:
  -highlight string
    	RGB color to highlight visible scratchpad window name (default "fff")
  -op string
    	operation on scratchpad: detect - to check if it exists, toggle - to toggle scratchpad (default "detect")
  -s string
    	scratchpad name
  -truncate int
    	truncate window title. 0 - to disable. (default 20)
```