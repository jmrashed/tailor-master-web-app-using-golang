package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

// main is the entry point for the application.
func main() {
	// start the web server
	serveWeb()
}

// struc to pass into the templage
type defaultContext struct {
	Title       string
	ErrorMsgs   string
	SuccessMsgs string
}

var themeName = getThemeName()          //late we will collect this value from config file
var staticPages = populateStaticPages() //custom function to colect all the available pages under pages folder

// serveWeb starts the web server and registers the necessary handlers.
func serveWeb() {

	// create a new Gorilla router for handling dynamic URL routes
	gorillaRoute := mux.NewRouter()

	// register the serveContent function to handle the root URL and dynamic page aliases
	gorillaRoute.HandleFunc("/", serveContent)
	gorillaRoute.HandleFunc("/{page_alias}", serveContent)

	// register the serveResource function to handle requests for image, CSS, and JavaScript files
	http.HandleFunc("/img/", serveResource)
	http.HandleFunc("/css/", serveResource)
	http.HandleFunc("/js/", serveResource)

	// register the Gorilla router as the handler for all other routes
	http.Handle("/", gorillaRoute)

	// start the HTTP server and listen for incoming requests on port 8080
	http.ListenAndServe(":8080", nil)
	fmt.Println("Server listening on http://localhost:8080/")

}

// serveContent serves the requested content based on the URL parameter "page_alias".
func serveContent(w http.ResponseWriter, r *http.Request) {
	// retrieve the page_alias parameter from the URL
	urlParams := mux.Vars(r)
	page_alias := urlParams["page_alias"]
	// set the default page_alias to "home" if it's empty
	if page_alias == "" {
		page_alias = "home"
	}
	// lookup the static page template by concatenating the page_alias with the ".html" extension
	staticPage := staticPages.Lookup(page_alias + ".html")
	log.Println("page ", staticPage)

	// if the static page template is not found, use the 404 template instead and return a 404 error code
	if staticPage == nil {
		staticPage = staticPages.Lookup("404.html")
		w.WriteHeader(404)
	}

	// set the context variables to be passed into the template
	context := defaultContext{}
	context.Title = page_alias
	context.ErrorMsgs = ""
	context.SuccessMsgs = ""
	log.Println(context)

	// execute the static page template with the context variables and write the output to the response writer
	err := staticPage.Execute(w, context)
	if err != nil {
		log.Println(err)
	}

}
func getThemeName() string {
	return "bs5"
}

// --------------------------------------------------------------
// populateStaticPages retrieves all the files under the given folder and its subsequent folders to populate the static pages for the website
func populateStaticPages() *template.Template {
	result := template.New("templates")
	templatePaths := new([]string)

	// set the base path for the page templates
	basePath := "pages"
	// open the folder containing the page templates
	templateFolder, _ := os.Open(basePath)
	defer templateFolder.Close()
	// get a list of all files in the folder
	templatePathsRaw, _ := templateFolder.Readdir(-1)
	// iterate over each file and add its path to the list of template paths
	for _, pathInfo := range templatePathsRaw {
		log.Println(pathInfo.Name())
		*templatePaths = append(*templatePaths, basePath+"/"+pathInfo.Name())
	}

	// set the base path for the theme templates
	basePath = "themes/" + themeName
	// open the folder containing the theme templates
	templateFolder, _ = os.Open(basePath)
	defer templateFolder.Close()
	// get a list of all files in the folder
	templatePathsRaw, _ = templateFolder.Readdir(-1)
	// iterate over each file and add its path to the list of template paths
	for _, pathInfo := range templatePathsRaw {
		log.Println(pathInfo.Name())
		*templatePaths = append(*templatePaths, basePath+"/"+pathInfo.Name())
	}

	// parse all the template files and return the result
	result.ParseFiles(*templatePaths...)
	return result
}

//--------------------------------------------------------------

//--------------------------------------------------------------

// serveResource serves the requested resource of types js, img, css files.
func serveResource(w http.ResponseWriter, req *http.Request) {
	// build the file path by concatenating the "public" folder, the theme name, and the URL path
	path := "public/" + themeName + req.URL.Path
	var contentType string
	// determine the content type based on the file extension
	if strings.HasSuffix(path, ".css") {
		contentType = "text/css; charset=utf-8"
	} else if strings.HasSuffix(path, ".png") {
		contentType = "img/png; charset=utf-8"
	} else if strings.HasSuffix(path, ".jpg") {
		contentType = "img/jpg; charset=utf-8"
	} else if strings.HasSuffix(path, ".js") {
		contentType = "application/javascript; charset=utf-8"
	} else {
		contentType = "text/plain; charset=utf-8"
	}

	// log the path of the requested file
	log.Println(path)
	// open the file
	f, err := os.Open(path)
	if err == nil {
		defer f.Close()
		// set the content type header
		w.Header().Add("Content-Type", contentType)
		// create a buffered reader for the file
		br := bufio.NewReader(f)
		// write the contents of the file to the response writer
		br.WriteTo(w)
	} else {
		// return a 404 error if the file is not found
		w.WriteHeader(404)
	}
}

//--------------------------------------------------------------
