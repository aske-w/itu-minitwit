# Steps
- Install python3 and pip
- Install dependencies: `pip install -r requirements.txt`
- Install GCC
 - Install SQLite: `sudo apt install libsqlite3-dev`
 - Install sqlitebrowser: `sudo apt install sqlitebrowser`
 - Compile **flag_tool.c** with: `gcc flag_tool.c -l sqlite3 -o flag_tool`


## Download DB
`scp student@165.227.167.15:/tmp/minitwit.db ~/Desktop/itu-minitwit`
password: uiuiui