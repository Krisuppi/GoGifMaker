# Instructions

- Install ffmpeg
- Put the app in a folder with images (jpg, jpeg, png supported) that are sized same. May cause unwanted result if different but will attempt to do so with a warning message.
- Launch it (go run or build yourself)
- Enjoy resulted gif

# Configure
- Default values: png files, 30fps, shows log
- Edit by creating config.txt in same folder as the app
- 1st line defines filetype without dot at start e.g "png"
- 2nd line defines fps
- 3rd line defines debug mode if any text
- 4+ ignored

# How it works
- Ffmpeg is a wonderful tool but defining input range for gifs or videos can be a pain when files are not named in a specific pattern.
- This app temporarily renames all files in folder matching filetype and runs ffmpeg on them. Then renames them back. Sets the ordering for gif with default cmp. Compares ignoring case sensitivity.
- In order to deal with transparency gif is done in 3 stages, creating tiles from which palette will be created and using palette ffmpeg converts images to gif. Finally temp files for tile and palette are deleted.

# TODOs
- Tests with turning script-like structure to go-like
- Resize result if mismatches are found
- Invalid filetype/fps to warn that it will use png/30fps.
- Check for dirs (maybe)
- More filetypes for input and output (maybe)
- Dig into log-levels. Maybe could remove explicit if debug then log statements with this (maybe)

# Done
- Png and jpg/jpeg support to gif
- Adjustable fps
- Hide output with debug option off
- Transparency preservation with png using tile and palette generation
- Validate config input. Validation will cause defaults to be used
- Debug akshually useful. Prints resulting ffmpeg cmds and used config. On error prints meaningful error. 

