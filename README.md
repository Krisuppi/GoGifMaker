# Instructions

- Install ffmpeg
- Put the app in a folder with images
- Launch it
- Enjoy gif

# Configure
- Default values: .png files, 30fps, shows log
- Edit by creating config.txt in same folder as the app
- 1st line defines filetype with dot at start e.g ".png"
- 2nd line defines fps
- 3rd line defines debug mode if any text but whitespace is set

# How it works
- ffmpeg is a wonderful tool but defining input range for gifs or videos can be a pain when files are not named in a specific pattern.
- This app temporarily renames all files in folder matching filetype and runs ffmpeg on them. Then renames them back. Sets the ordering for gif with default cmp.Compare ignoring case sensitivity.

# TODOs
- maybe mp4 if waifu wants
- sanitize input in case this gets on a server
