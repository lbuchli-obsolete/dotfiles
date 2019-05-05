package main

/*
##############################################################
# Section: Imports
##############################################################
*/

import (
	"bufio"
	"crypto/md5"
	"encoding/binary"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

/*
##############################################################
# Section: Constants & Fields
##############################################################
*/

var display_size Vector

var background_color uint32

// The path to the folder containing images of the current desktops
const DESKTOP_IMAGES_PATH = "~/.config/sway/dimgs/"

/*
##############################################################
# Section: Basic Types and functions
##############################################################
*/

type Vector struct {
	x int32
	y int32
}

type FractionVector struct {
	x float32
	y float32
}

// This wants to be an enum
type Align int

const (
	LEFT   Align = 0
	TOP    Align = 0
	CENTER Align = 1
	RIGHT  Align = 2
	BOTTOM Align = 2
)

// text sizes
const (
	TITLE     int = 128
	SUBTITLE  int = 64
	HEADER    int = 32
	SUBHEADER int = 24
	TEXT      int = 16
	SUBTEXT   int = 14
)

func UInt32ToColor(ui uint32) (color sdl.Color) {
	bytes := (*[4]byte)(unsafe.Pointer(&ui))[:]
	return sdl.Color{R: bytes[2], G: bytes[1], B: bytes[0], A: bytes[3]}
}

func ImgTosurface(img image.Image) (surface *sdl.Surface, err error) {
	// Credit to https://github.com/veandco/go-sdl2/issues/116#issuecomment-96056082
	rgba := image.NewRGBA(img.Bounds())
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	s, err := sdl.CreateRGBSurface(0, int32(w), int32(h), 32, 0, 0, 0, 0)
	if err != nil {
		return s, err
	}
	rgba.Pix = s.Pixels()

	for y := 0; y < h; y += 1 {
		for x := 0; x < w; x += 1 {
			c := img.At(x, y)
			rgba.Set(x, y, c)
		}
	}

	return s, nil
}

func resizeSurface(surf *sdl.Surface, newsize Vector) (resizedsurf *sdl.Surface, err error) {
	resizedsurf, err = sdl.CreateRGBSurface(0, newsize.x, newsize.y, 32, 0, 0, 0, 0)
	if err != nil {
		return resizedsurf, err
	}

	rgba := image.NewRGBA(image.Rect(0, 0, int(newsize.x), int(newsize.y)))
	rgba.Pix = resizedsurf.Pixels()

	pixels := surf.Pixels()

	// the factor from newscale to oldscale
	scalex := float32(surf.W) / float32(newsize.x)
	scaley := float32(surf.H) / float32(newsize.y)

	for y := int32(0); y < newsize.y; y += 1 {
		for x := int32(0); x < newsize.x; x += 1 {
			// the byte index of the color currently looking at
			//		*					line 											  *  			 pixel in line (x)					*
			index := int32((float32(surf.BytesPerPixel()) * (float32(surf.W*y) * scaley)) + (float32(int32(surf.BytesPerPixel())*x) * scalex))

			// set the pixel using magic
			rgba.Set(int(x), int(y), UInt32ToColor(binary.BigEndian.Uint32(pixels[index:index+4])))
		}
	}

	return resizedsurf, err
}

/*
###############################################################
# Section: Initialization
###############################################################
*/

func Initialize() {
	// Initialize sdl.sdl and sdl.ttf
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ttf.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get the display size
	bounds, err := sdl.GetDisplayBounds(0)
	if err != nil {
		fmt.Println(err)
		return
	}

	display_size = Vector{bounds.W, bounds.H}
}

/*
###############################################################
# Section: Item & Container
###############################################################
*/

// Every type with a position and scale is considered an item.
type Item interface {
	Draw(*sdl.Surface) error
	GetPosition() Vector
	SetPosition(Vector)
	GetSize() Vector
	SetSize(Vector)
}

// This is the first (and the most important) item.
// It is used to group other items.
type Container struct {
	position Vector
	size     Vector
	items    map[string]Item
}

// Move the item to a pixel position
func (cont *Container) MoveItem(item string, pos Vector) {
	cont.items[item].SetPosition(pos)
}

// Move the item to a fraction of the parent container size
func (cont *Container) MoveItemToFraction(item string, pos FractionVector) {
	cont.items[item].SetPosition(Vector{int32(pos.x * float32(cont.size.x)), int32(pos.y * float32(cont.size.y))})
}

// Resize an Item to a specific pixel size
func (cont *Container) ResizeItem(item string, size Vector) {
	cont.items[item].SetSize(size)
}

// Resize an Item to a fraction of the parent container size
func (cont *Container) ResizeItemToFraction(item string, size FractionVector) {
	cont.items[item].SetSize(Vector{int32(size.x * float32(cont.size.x)), int32(size.y * float32(cont.size.y))})
}

// draw a container
// The container will let each item draw onto its own surface and then draw that onto the main surface
func (cont *Container) Draw(surf *sdl.Surface) (err error) {

	// let each item draw onto the surface
	for _, val := range cont.items {
		isurface, err := sdl.CreateRGBSurface(0, val.GetSize().x, val.GetSize().y, 32, 0, 0, 0, 0)
		if err != nil {
			return err
		}

		// also apply background color
		isurface.FillRect(nil, background_color)

		err = val.Draw(isurface)
		if err != nil {
			return err
		}

		pos := val.GetPosition()
		size := val.GetSize()

		// draw the item surface onto the container surface
		src_rect := sdl.Rect{X: 0, Y: 0, W: size.x, H: size.y}
		dst_rect := sdl.Rect{X: pos.x, Y: pos.y, W: pos.x + size.x, H: pos.y + size.y}
		isurface.Blit(&src_rect, surf, &dst_rect)

		isurface.Free()
	}

	return nil
}

// Add an item to the container
func (cont *Container) AddItem(name string, item Item) {
	cont.items[name] = item
}

// Get an item from the container
func (cont *Container) GetItem(name string) (item Item) {
	return cont.items[name]
}

// Getters and setters

func (cont *Container) GetPosition() (position Vector) {
	return cont.position
}

func (cont *Container) SetPosition(position Vector) {
	cont.position = position
}

func (cont *Container) GetSize() (size Vector) {
	return cont.size
}

func (cont *Container) SetSize(size Vector) {
	cont.size = size
}

/*
####################################################################
# Section: Basic item types
####################################################################
*/

/*
########################
# Subsection: Label
########################
*/

type Label struct {
	position Vector
	size     Vector
	text     string
	textsize int
	valign   Align
	halign   Align
	color    uint32
	bgcolor  uint32
	bold     bool
}

// Draw the item onto the parent surface
func (label *Label) Draw(surf *sdl.Surface) (err error) {

	var font *ttf.Font

	// load font
	if label.bold {
		font, err = ttf.OpenFont("/usr/share/fonts/TTF/DejaVuSans-Bold.ttf", label.textsize)
	} else {
		font, err = ttf.OpenFont("/usr/share/fonts/TTF/DejaVuSans.ttf", label.textsize)
	}
	if err != nil {
		return err
	}

	// Render text to surface
	text_surface, err := font.RenderUTF8Shaded(label.text, UInt32ToColor(label.color), UInt32ToColor(label.bgcolor))
	if err != nil {
		return err
	}
	defer text_surface.Free()

	// Calculate vertical and horizontal position on surface
	var coordinate_x int32
	var coordinate_y int32

	switch label.halign {
	case LEFT:
		coordinate_x = 0
	case CENTER:
		coordinate_x = (int32(label.size.x) - text_surface.W) / 2
	case RIGHT:
		coordinate_x = int32(label.size.x) - text_surface.W
	}

	switch label.valign {
	case TOP:
		coordinate_y = 0
	case CENTER:
		coordinate_y = (int32(label.size.y) - text_surface.H) / 2
	case BOTTOM:
		coordinate_y = int32(label.size.y) - text_surface.H
	}

	dst_rect := sdl.Rect{X: coordinate_x, Y: coordinate_y, W: coordinate_x + text_surface.W, H: coordinate_y + text_surface.H}

	surf.FillRect(nil, background_color)

	// Draw onto final surface (Text aligned)
	text_surface.Blit(&sdl.Rect{X: 0, Y: 0, W: text_surface.W, H: text_surface.H}, surf, &dst_rect)

	return nil
}

// Getters and setters

func (label *Label) GetPosition() (position Vector) {
	return label.position
}

func (label *Label) SetPosition(newposition Vector) {
	label.position = newposition
}

func (label *Label) GetSize() (size Vector) {
	return label.size
}

func (label *Label) SetSize(newsize Vector) {
	label.size = newsize
}

/*
########################
# Subsection: Texture
########################
*/

type Texture struct {
	position Vector
	size     Vector
	texture  *sdl.Surface
}

// Draw the item onto the parent surface
func (tex *Texture) Draw(surf *sdl.Surface) (err error) {
	src_rect := sdl.Rect{X: 0, Y: 0, W: tex.size.x, H: tex.size.y}
	dst_rect := sdl.Rect{X: 0, Y: 0, W: tex.size.x, H: tex.size.y}
	tex.texture.Blit(&src_rect, surf, &dst_rect)

	return nil
}

// Getters and setters

func (tex *Texture) GetPosition() (position Vector) {
	return tex.position
}

func (tex *Texture) SetPosition(position Vector) {
	tex.position = position
}

func (tex *Texture) GetSize() (size Vector) {
	return tex.size
}

func (tex *Texture) SetSize(size Vector) {
	tex.size = size
}

/*
########################
# Subsection: Unicolor
########################
*/

type Unicolor struct {
	position Vector
	size     Vector
	color    uint32
}

// Draw the item onto the parent surface
func (unic *Unicolor) Draw(surf *sdl.Surface) (err error) {
	rect := sdl.Rect{X: 0, Y: 0, W: unic.size.x, H: unic.size.y}
	return surf.FillRect(&rect, unic.color)
}

// Getters and setters

func (unic *Unicolor) GetPosition() (position Vector) {
	return unic.position
}

func (unic *Unicolor) SetPosition(position Vector) {
	unic.position = position
}

func (unic *Unicolor) GetSize() (size Vector) {
	return unic.size
}

func (unic *Unicolor) SetSize(size Vector) {
	unic.size = size
}

/*
##############################################################
# Section: Window
##############################################################
*/

type WindowHandler interface {
	Init(*Container, *bool)
	Update()
	HandleEvent(sdl.Event)
}

func CreateWindow(position Vector, size Vector, bgcolor uint32, handler WindowHandler) (err error) {
	// This variable will will determine wether the window is running or not
	running := true

	// the main container
	cont := Container{Vector{0, 0}, size, make(map[string]Item)}

	// create an sdl window for the window struct instance
	window, err := sdl.CreateWindow("Sidebar", position.x, position.y,
		cont.size.x, cont.size.y, sdl.WINDOW_POPUP_MENU)
	if err != nil {
		return err
	}

	surface, err := window.GetSurface()
	if err != nil {
		return err
	}
	defer window.Destroy()
	defer sdl.Quit()

	// Set the background color
	surface.FillRect(nil, bgcolor)
	background_color = bgcolor

	// Initialize the handler
	handler.Init(&cont, &running)

	// The main loop
	for running {

		// Quit the program in case of exit event
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("Exit signal received. Quitting...")
				running = false
				break
			default:
				handler.HandleEvent(event)
			}
		}

		handler.Update()
		cont.Draw(surface)
		window.UpdateSurface()
	}

	return nil
}

/*
######################################################################################################
######################################################################################################
## Chapter: Windows                                                                                 ##
######################################################################################################
######################################################################################################
*/

/*
####################################################################
# Section: Main
####################################################################
*/

const DEF_BG_COLOR uint32 = 0x10171e
const WHITE_COLOR uint32 = 0xffffff

const SCREEN_FRACTION = 4

func main() {
	// initialize packages
	Initialize()

	// 0 is the program itself, 1 is the first argument
	var arg string
	if len(os.Args) > 1 {
		arg = os.Args[1]
	} else {
		arg = "notspecified"
	}

	var handler WindowHandler

	// Determine the type of window
	switch arg {
	case "power":
		handler = &PowerWindowHandler{}
	case "run":
		handler = &RunWindowHandler{}
	case "desktop":
		handler = &DesktopWindowHandler{}
	default:
		handler = &RunWindowHandler{}
	}

	CreateWindow(Vector{0, 0}, Vector{display_size.x / SCREEN_FRACTION, display_size.y}, DEF_BG_COLOR, handler)

	//path, err := getDataFilePath()
	//if err != nil {
	//panic(err)
	//}

	//fmt.Println(incrementDataFileEntry(path, "/usr/share/applications/firefox.desktop"))
}

/*
#################################################################
# Section: Power
#################################################################
*/

type PowerWindowHandler struct {
	cont *Container
	exit *bool
}

func (pwh *PowerWindowHandler) Init(c *Container, e *bool) {
	pwh.cont = c
	pwh.exit = e

	pwh.cont.AddItem("title", &Label{
		position: Vector{0, 0},
		size:     Vector{0, 0}, // will be resized later
		text:     "Power",
		textsize: 128,
		valign:   CENTER,
		halign:   CENTER,
		color:    WHITE_COLOR,
		bgcolor:  DEF_BG_COLOR,
		bold:     false,
	})

	pwh.cont.ResizeItemToFraction("title", FractionVector{1.0, 0.1})
}

func (pwh *PowerWindowHandler) Update() {
	return
}

func (pwh *PowerWindowHandler) HandleEvent(event sdl.Event) {
	return
}

/*
###################################################################
# Section: Run
###################################################################
*/

type RunWindowHandler struct {
	cont *Container
	exit *bool
}

func (rwh *RunWindowHandler) Init(c *Container, e *bool) {
	rwh.cont = c
	rwh.exit = e

	rwh.cont.AddItem("title", &Label{
		position: Vector{0, 0},
		size:     Vector{0, 0}, // will be resized later
		text:     "Run",
		textsize: 128,
		valign:   CENTER,
		halign:   CENTER,
		color:    WHITE_COLOR,
		bgcolor:  DEF_BG_COLOR,
		bold:     false,
	})

	rwh.cont.ResizeItemToFraction("title", FractionVector{1.0, 0.1})

	program, err := rwh.getProgramInfoCont("/usr/share/applications/nvim.desktop")
	if err != nil {
		fmt.Println(err)
	}

	rwh.cont.AddItem("program", program)
	rwh.cont.MoveItemToFraction("program", FractionVector{0, 0.1})

	program2, err := rwh.getProgramInfoCont("/usr/share/applications/termite.desktop")
	if err != nil {
		fmt.Println(err)
	}

	rwh.cont.AddItem("program2", program2)
	rwh.cont.MoveItemToFraction("program2", FractionVector{0, 0.1 + (float32(1) / float32(16))})
}

func (rwh *RunWindowHandler) Update() {
	return
}

func (rwh *RunWindowHandler) HandleEvent(event sdl.Event) {
	switch ty := event.(type) {
	// Does not capture keyboard event (possibly because of window type)
	// TODO find some other method to get keyboard input
	case *sdl.KeyboardEvent:
		fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
			ty.Timestamp, ty.Type, ty.Keysym.Sym, ty.Keysym.Mod, ty.State, ty.Repeat)
	}
}

// Gets a container containing info about a .desktop file
func (rwh *RunWindowHandler) getProgramInfoCont(path string) (cont *Container, err error) {
	cont = &Container{Vector{0, 0}, Vector{rwh.cont.size.x, rwh.cont.size.y / 16}, make(map[string]Item)}

	info, err := parseDesktopFile(path)
	if err != nil {
		return cont, err
	}

	// separator bar
	cont.AddItem("bar", &Unicolor{
		position: Vector{0, 0},
		size:     Vector{cont.size.x, 4},
		color:    0x20272e,
	})

	iconsize := rwh.cont.size.y / 18

	var icon *sdl.Surface

	default_icon := func() {
		icon, err = sdl.CreateRGBSurface(0, iconsize, iconsize, 32, 0, 0, 0, 0)
		icon.FillRect(nil, DEF_BG_COLOR)
	}

	homepath, err := getHomePath()
	if err != nil {
		return cont, err
	}

	// open the icon file. If anything fails,
	// use an empty surface instead
	// ugly code incoming
	// TODO find way to get icon path (icons are in subdirectories of dirs below)
	iconfile, err := os.Open("/usr/share/icons/" + info["Icon"])
	if err != nil {
		iconfile, err := os.Open(homepath + "/.local/share/icons/" + info["Icon"])
		if err != nil {
			default_icon()
		} else {
			img, _, err := image.Decode(iconfile)
			if err != nil {
				default_icon()
			} else {
				iconsurf, err := ImgTosurface(img)
				if err != nil {
					default_icon()
				} else {
					icon, err = resizeSurface(iconsurf, Vector{iconsize, iconsize})
					if err != nil {
						default_icon()
					}
				}
			}
		}
	}

	// if iconfile was initialized, close it again
	if iconfile != nil {
		err = iconfile.Close()
		if err != nil {
			return cont, err
		}
	}

	// icon of the program
	cont.AddItem("icon", &Texture{
		position: Vector{0, 8},
		size:     Vector{iconsize, iconsize},
		texture:  icon,
	})

	// name of the program
	cont.AddItem("title", &Label{
		position: Vector{iconsize + 8, 8},
		size:     Vector{0, 0}, // will be resized later
		text:     info["Name"],
		textsize: HEADER,
		valign:   TOP,
		halign:   LEFT,
		color:    WHITE_COLOR,
		bgcolor:  DEF_BG_COLOR,
		bold:     false,
	})
	cont.ResizeItemToFraction("title", FractionVector{float32(16) / float32(18), 0.5})

	// description of the program
	cont.AddItem("description", &Label{
		position: Vector{iconsize + 8, (cont.size.y / 2) + 8},
		size:     Vector{0, 0}, // will be resized later
		text:     info["Comment"],
		textsize: SUBHEADER,
		valign:   BOTTOM,
		halign:   LEFT,
		color:    0xa0a1a7,
		bgcolor:  DEF_BG_COLOR,
		bold:     false,
	})
	cont.ResizeItemToFraction("description", FractionVector{float32(16) / float32(18), 0.2})

	return cont, err
}

/*
####################################################################
# Section: Desktop
####################################################################
*/

type DesktopWindowHandler struct {
	cont *Container
	exit *bool
}

func (dwh *DesktopWindowHandler) Init(c *Container, e *bool) {
	dwh.cont = c
	dwh.exit = e

	dwh.cont.AddItem("title", &Label{
		position: Vector{0, 0},
		size:     Vector{0, 0}, // will be resized later
		text:     "Desktops",
		textsize: 128,
		valign:   CENTER,
		halign:   CENTER,
		color:    WHITE_COLOR,
		bgcolor:  DEF_BG_COLOR,
		bold:     false,
	})
	dwh.cont.ResizeItemToFraction("title", FractionVector{1.0, 0.1})

	// add desktop images
	for i := 1; i <= 6; i++ {
		var img_surface *sdl.Surface
		defer img_surface.Free()

		// get the surface
		file, err := os.Open(DESKTOP_IMAGES_PATH + strconv.Itoa(i) + ".png")
		if err != nil {
			// assuming there was no image or image is corrupted; display empty space
			img_surface, err = GetEmptyDesktop()
			if err != nil {
				panic(err)
			}
		} else {
			img, _, err := image.Decode(file)
			if err != nil {
				img_surface, err = GetEmptyDesktop()
				if err != nil {
					panic(err)
				}
			} else {
				img_surface, err = ImgTosurface(img)
				if err != nil {
					continue
				}
			}

		}

		abs_pos := Vector{int32((i - 1) % 2), int32(math.Ceil(float64(i)/2.0)) - 1}

		resized_surface, err := resizeSurface(img_surface, Vector{display_size.x / 12, display_size.y / 12})
		if err != nil {
			panic(err)
		}

		desktop_cont := &Container{
			position: Vector{(display_size.x / 8) * abs_pos.x,
				int32(float32(display_size.y)*0.1) + (int32((float32(display_size.y)*0.9)/6) * abs_pos.y)},
			size: Vector{display_size.x / 8, int32(float32(display_size.y)*0.9) / 6},
			items: map[string]Item{
				"number": &Label{
					position: Vector{0, 0},
					size:     Vector{0, 0}, // will be resized
					text:     strconv.Itoa(i),
					textsize: 64,
					valign:   TOP,
					halign:   CENTER,
					color:    WHITE_COLOR,
					bgcolor:  DEF_BG_COLOR,
					bold:     true,
				},
				"image": &Texture{
					position: Vector{0, 0}, // will be repositioned
					size:     Vector{0, 0}, // will be resized
					texture:  resized_surface,
				},
			},
		}

		// do resizing and repositioning of above mentioned items
		desktop_cont.ResizeItemToFraction("number", FractionVector{0.2, 1})

		desktop_cont.MoveItemToFraction("image", FractionVector{0.2, 0})
		desktop_cont.ResizeItemToFraction("image", FractionVector{0.8, 1})

		// add the container to the parent container
		dwh.cont.AddItem("desktop-"+strconv.Itoa(i), desktop_cont)
	}
}

// Gets you a surface the size of the current desktop, uniform colored
func GetEmptyDesktop() (desktop *sdl.Surface, err error) {
	desktop, err = sdl.CreateRGBSurface(0, display_size.x, display_size.y, 32, 0, 0, 0, 0)
	if err != nil {
		return desktop, err
	}
	desktop.FillRect(nil, uint32(0x20272e))
	return desktop, err
}

func (dwh *DesktopWindowHandler) Update() {
	return
}

func (dwh *DesktopWindowHandler) HandleEvent(event sdl.Event) {
	return
}

/*
######################################################################################################
######################################################################################################
## Chapter: Desktop File Parser                                                                     ##
######################################################################################################
######################################################################################################
*/

var desktop_file_paths = []string{"/usr/local/share/applications/", "/usr/share/applications/", "~/.local/share/applications/"}

func parseDesktopFile(file string) (entries map[string]string, err error) {
	entries = make(map[string]string)

	lines, err := readFile(file)
	if err != nil {
		return entries, err
	}

	for _, line := range lines {
		matched, err := regexp.MatchString("([a-zA-Z]+)=(.*)", line)
		if err != nil {
			return entries, err
		}

		// if the line is a valid entry
		if matched {

			// Searches for a sequence of letters with length >= 1
			key_regex, err := regexp.Compile("^[a-zA-Z]+")
			if err != nil {
				return entries, err
			}

			// Searches for a sequence of characters with = as a prefix
			value_regex, err := regexp.Compile("=(.*)")
			if err != nil {
				return entries, err
			}

			// make entry but leave out the '='
			entries[key_regex.FindString(line)] = value_regex.FindString(line)[1:]

		}
	}

	return entries, nil
}

func readFile(filepath string) (lines []string, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return lines, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return lines, err
	}

	return lines, nil
}

func getFileHashes(dir_paths []string) (file_hashes map[string]string, err error) {
	file_hashes = make(map[string]string)

	for _, dir_path := range dir_paths {

		filepath.Walk(dir_path, func(path string, info os.FileInfo, err error) error {

			// for some reason this function gets called sometimes with no fileinfo
			// so we have to check for that
			if info == nil {
				return nil
			}

			// do not process folders
			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}

			// calculate the md5 hash of file
			hash := md5.New()
			if _, err := io.Copy(hash, file); err != nil {
				return err
			}

			// Close the file
			err = file.Close()
			if err != nil {
				return err
			}

			// assign the file path to the hash
			file_hashes[hex.EncodeToString(hash.Sum(nil))] = path

			return nil
		})
	}

	return file_hashes, nil
}

// Gets all the files in a list of directories as keys in a map
func getFiles(dir_paths []string) (file_paths map[string]struct{}, err error) {
	file_paths = make(map[string]struct{})

	for _, dir_path := range dir_paths {
		filepath.Walk(dir_path, func(path string, info os.FileInfo, err error) error {

			// sometimes this function gets called with no fileinfo
			// so we have to ignore these calls
			if info == nil {
				return err
			}

			// also ignore folders
			if info.IsDir() {
				return err
			}

			// assign an empty struct to the key
			file_paths[path] = struct{}{}

			return err
		})
	}

	return file_paths, err
}

func checkEntryMatch(key string, value string, parsedfile map[string]string) (matches bool) {
	val, ok := parsedfile[key]

	if ok {
		// Check match ignoring case
		if strings.EqualFold(val, value) {
			return true
		}
	}

	return false
}

func launchDesktopFile(path string) (err error) {
	return exec.Command("gtk-launch", path).Run()
}

/*
######################################################################################################
######################################################################################################
## Chapter: CSV Data Processing																		##
######################################################################################################
######################################################################################################
*/

// gets the path to the data file
func getDataFilePath() (path string, err error) {
	// The value regexp
	value_regex, err := regexp.Compile("=(.*)")
	if err != nil {
		return "", err
	}

	// the variable regexp
	var_regex, err := regexp.Compile("HOME=(.*)")
	if err != nil {
		return "", err
	}

	for _, line := range os.Environ() {
		// if its the entry we search for ($HOME)
		if var_regex.MatchString(line) {
			// return the value of the env variable + our config file location
			return value_regex.FindString(line)[1:] + "/sway/sidebar/data.csv", err
		}
	}

	return "", errors.New("Enviroment variable $HOME not found")
}

func getHomePath() (path string, err error) {
	// The value regexp
	value_regex, err := regexp.Compile("=(.*)")
	if err != nil {
		return "", err
	}

	// the variable regexp
	var_regex, err := regexp.Compile("HOME=(.*)")
	if err != nil {
		return "", err
	}

	for _, line := range os.Environ() {
		// if its the entry we search for ($HOME)
		if var_regex.MatchString(line) {
			// return the value of the env variable + our config file location
			return value_regex.FindString(line)[1:], err
		}
	}

	return "", errors.New("Enviroment variable $HOME not found")
}

// Reads out the data csv saved under CONFIG_FILE_PATH
func readDataFile(path string) (data map[string][]string, err error) {
	data = make(map[string][]string)

	// open or create file in readonly mode
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return data, err
	}

	reader := csv.NewReader(file)

	raw_data, err := reader.ReadAll()
	if err != nil {
		return data, err
	}

	// make the first entry of each line the line's key
	for _, line := range raw_data {
		data[line[0]] = line[1:]
	}

	return data, err
}

// overwrites the current datafile by replacing it
func writeDataFile(path string, data map[string][]string) (err error) {
	// delete the old data file
	deleteDataFile(path)

	// open or create file in writeonly mode
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)

	var records [][]string

	// reformat data to two dimensional array
	for key, val := range data {
		records = append(records, append([]string{key}, val...))
	}

	return writer.WriteAll(records)
}

// increments the use counter of a entry in the data file
func incrementDataFileEntry(path string, key string) (err error) {
	// get the data in a validated form
	data, err := getValidatedData(path)
	if err != nil {
		return err
	}

	i, err := strconv.Atoi(data[key][0])
	if err != nil {
		return err
	}

	data[key] = []string{strconv.Itoa(i + 1)}

	return writeDataFile(path, data)
}

// removes the current data file.
func deleteDataFile(path string) (err error) {
	return os.Remove(path)
}

// Same as readDataFile(), but validates entries first
func getValidatedData(path string) (data map[string][]string, err error) {
	data, err = readDataFile(path)
	if err != nil {
		return data, err
	}

	files, err := getFiles(desktop_file_paths)
	if err != nil {
		return data, err
	}

	// add new desktop files
	for path, _ := range files {
		// if the desktop file is missing in data
		if _, ok := data[path]; !ok {
			data[path] = []string{"0"}
		}
	}

	// remove no longer existent desktop files
	for path, _ := range data {
		// if there is no such file
		if _, ok := files[path]; !ok {
			delete(data, path)
		}
	}

	// return validated data variable
	return data, err
}
