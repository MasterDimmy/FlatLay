This Golang + HTML + JavaScript program shows solving of Flat Lay problem.
So it can create something like this: http://theexposure.co/wp-content/uploads/2015/07/Photo-27-06-2015-3-02-27-pm.jpg

Windows run: collages.exe
Linux, iOS run: go run 

Starts web-server at 7070 port.
Go to http://localhost:7070
Change the browser size - look at the result in main window.

How its made:
Program create_db.exe scans folder with console command
create_db.exe <path> and created the database of pictures "database.json" with array of objects

type TImage struct {
    Name   string //picrure name
    Path   string //path to picture
    Width  int    //
    Height int    //
    Weight int    //not used now (for top and bottom layout of selected pictures)
    Group  int    //group of pictures (cars, clothe's style, etc.)
}

Folders & files:
	looks - keeps pictures
	static - web-server
   index.html - demo page
   onload.js - working script

On run:
Script gets total groups count at "/get_total_groups".
Then script asks server to solve Flat Lay at given borders and get result at "/get_field?group=G&width=W&height=H"

Server's answer: JSON array of pictures and they offsets for Flat Lay.

Current solution has timeout for solving with too many pictures: 1 sec. It takes best variant it could find on that time.
If you set the group number, it returns best solution.

All server's output is in JSON with simple protocol.

Sorry for English...