call plug#begin('~/.config/nvim/bundle')
" List of all plugins to use
Plug 'Shougo/deoplete.nvim', {'do': ':UpdateRemotePlugins'} " Autocompletion
Plug 'scrooloose/nerdtree' " File overview
Plug 'scrooloose/syntastic' " Syntax highlighting
Plug 'tpope/vim-surround' " Surround words with parenthesis, etc...
Plug 'tpope/vim-fugitive' " Git integration for vim
Plug 'bling/vim-airline' " Infobar
Plug 'vim-airline/vim-airline-themes'
Plug 'scrooloose/nerdcommenter' " Comment out code
Plug 'jiangmiao/auto-pairs' " Adds parenthesis pair instead of just ( when typing
Plug 'zchee/deoplete-go', {'do': 'make'} " Auto-completion for golang
Plug 'zchee/deoplete-jedi' " Auto-completion for python
Plug 'majutsushi/tagbar' " overview of variables, etc..
Plug 'SirVer/ultisnips' " Snippets for autocompletion, ...
Plug 'honza/vim-snippets' " Actual snippets for ultisnips
Plug 'xuhdev/vim-latex-live-preview', { 'for': 'tex' } " live preview for LaTeX
Plug 'zchee/deoplete-clang' " Autocompletion for C/C++
Plug 'fatih/vim-go' " go-specific tools
Plug 'w0rp/ale' " code linting
Plug 'sebdah/vim-delve' " Integration of delve, a go debugger
Plug 'Yggdroot/indentLine' " Indentation guidelines
Plug 'airblade/vim-gitgutter' " Git diff in realtime
"Plug 'JamshedVesuna/vim-markdown-preview' " Markdown preview in browser
Plug 'mhinz/vim-startify' " Fancy start screen
Plug 'jceb/vim-orgmode' " Orgmode integration for vim
Plug 'othree/xml.vim' " XML tools
Plug 'lervag/vimtex' " LaTex
Plug 'vim-scripts/utl.vim' " Universal Text Linking (vim-orgmode dep)
Plug 'tpope/vim-repeat' " Better command repeatment (vim-orgmode dep)
Plug 'vim-scripts/taglist.vim' " Source code browsing (vim-orgmode dep)
Plug 'tpope/vim-speeddating' " Date tools (vim-orgmode dep)
Plug 'chrisbra/NrrwRgn' " Narrow Region feature from emacs (vim-orgmode dep)
Plug 'mattn/calendar-vim' " Calendar window inside vim (vim-orgmode dep)
Plug 'vim-scripts/SyntaxRange' " Syntax highlighting for code regions (vim-orgmode dep)
Plug 'calviken/vim-gdscript3' " Syntax highlighting and auto-completion for gdscript3
Plug 'Shougo/echodoc.vim' " Documentation in echo area
Plug 'stevearc/vim-arduino' " Arduino sketch tools & arduino uploading
Plug 'tbabej/taskwiki' " Project planning in vim
Plug 'blindFS/vim-taskwarrior' " TaskWarrior vim interface, taskwiki dep
Plug 'powerman/vim-plugin-AnsiEsc' " Ansi escape sequences, taskwiki dep
Plug 'vimwiki/vimwiki' " A personal wiki, taskwiki dep
Plug 'vhda/verilog_systemverilog.vim' " SystemVerilog plugin
Plug 'wannesm/wmgraphviz.vim' " GraphViz dot (compiling, viewing, ...)
Plug 'fidian/hexmode' " Hex editing
Plug 'tikhomirov/vim-glsl' " GLSL shader language syntax highlighting
Plug 'artur-shaik/vim-javacomplete2' " Java auto completion
Plug 'chenillen/jad.vim' " java decompiler
Plug 'neovimhaskell/haskell-vim' " Haskell syntax highlighting & indentation
Plug 'ctrlpvim/ctrlp.vim' " Fuzzy search
Plug 'vim-scripts/c.vim' " C programming language support
call plug#end()

" Initialization

"" Auto start NERD tree when opening a directory
"autocmd VimEnter * if argc() == 1 && isdirectory(argv()[0]) && !exists("s:std_in") | exe 'NERDTree' argv()[0] | wincmd p | ene | wincmd p | endif

"" Auto start NERD tree if no files are specified
"autocmd StdinReadPre * let s:std_in=1
"autocmd VimEnter * if argc() == 0 && !exists("s:std_in") | exe 'NERDTree' | endif

" Let quit work as expected if after entering :q the only window left open is NERD Tree itself
autocmd bufenter * if (winnr("$") == 1 && exists("b:NERDTree") && b:NERDTree.isTabTree()) | q | endif

" Basic variables
filetype plugin indent on
syntax on
let mapleader = "\<space>"
set incsearch
set nohlsearch
set smartcase
set tabstop=4
set softtabstop=0
set noexpandtab
set shiftwidth=4
set nowrap
set splitright
set splitbelow
set noshowmode
"set exrc
"set secure
let g:deoplete#sources#go#gocode_binary = $HOME.'/go/bin/gocode'

" No line numbers in text files or console
set number
autocmd FileType txt setlocal nonumber
autocmd TermOpen * setlocal nonumber

" Make split windows look more clean
set fillchars+=vert:\ 
hi VertSplit ctermfg=Gray

"#######Preferences########
" NERDTree
let NERDTreeQuitOnOpen = 1
nnoremap <leader>t :NERDTreeToggle<CR>
let NERDTreeAutoDeleteBuffer = 1 " Delete Buffer when file deleted
let NERDTreeMinimalUI = 1
let NERDTreeDirArrows = 1
let NERDTreeShowHidden=1 " Show hidden files
let NERDTreeIgnore=['\.DS_Store', '\~$', '\.swp'] " Ignore useless files

" Vim
" noremap Y 0y$
set hidden
set history=100
"autocmd BufWritePre * :%s/\s\+$//e " Delete whitespaces on saving
" Cancel searches with Escape
nnoremap <silent> <Esc> :nohlsearch<Bar>:echo<CR>
" Reopen previously opened file
nnoremap <Leader><Leader> :e#<CR>
" Highlight trailing whitespaces
autocmd ColorScheme * highlight ExtraWhitespace ctermbg=grey guibg=grey
" get rid of ~
"highlight EndOfBuffer ctermfg=bg ctermbg=NONE

" Autocompletion colors
highlight Pmenu ctermbg=0 ctermfg=15
highlight PmenuSel ctermbg=1 ctermfg=15
highlight PmenuSbar ctermbg=8 ctermfg=15

hi MatchParen cterm=bold ctermbg=None ctermfg=Red

" Other
nnoremap <Leader>b :TagbarToggle<CR>
let g:airline_theme='deus'
let g:airline_powerline_fonts = 1
" Make Tagbar-Highlighting look nicer
highlight TagbarHighlight ctermfg=yellow
set clipboard+=unnamedplus " Use the system clipboard

highlight Visual ctermbg=darkgrey cterm=bold

highlight SignColumn ctermbg=None

" Indentation guides
let g:indentLine_char = '|'
let g:indentLine_color_term = 239
let g:indentLine_fileType = ['go', 'xml', 'html', 'python', 'c', 'cpp', 'sh']

" Autocompletion
let g:deoplete#enable_at_startup = 1
" No autocompletion preview
set completeopt-=preview
" Allow saving of files as sudo when I forgot to start vim using sudo.
cmap w!! w !sudo tee > /dev/null %

" Disable deoplete when in multi cursor mode
function! Multiple_cursors_before()
    let b:deoplete_disable_auto_complete = 1
endfunction

function! Multiple_cursors_after()
    let b:deoplete_disable_auto_complete = 0
endfunction

" Faster switching of split window focus
"nnoremap <silent><Leader><Left> <c-w>h 
"nnoremap <silent><Leader><Right> <c-w>l 
"nnoremap <silent><Leader><Up> <c-w>k
"nnoremap <silent><Leader><Down> <c-w>j

nnoremap <silent><Leader>h <c-w>h 
nnoremap <silent><Leader>l <c-w>l 
nnoremap <silent><Leader>k <c-w>k
nnoremap <silent><Leader>j <c-w>j

" Disable arrow keys for training purposes
noremap <Up> <Nop>
noremap <Down> <Nop>
noremap <Left> <Nop>
noremap <Right> <Nop>

" ################# Running Programs ##################
autocmd FileType go nnoremap <buffer> <Leader>r :!go<Space>run<Space>%<Enter>
autocmd FileType go nnoremap <buffer> <Leader>R :DlvDebug<Enter>
autocmd FileType go nnoremap <buffer> <Leader><C-r> :DlvTest<Enter>
autocmd FileType go nnoremap <buffer> <Leader>g :GoDef<CR>

autocmd FileType arduino nnoremap <buffer> <leader>am :ArduinoVerify<CR>
autocmd FileType arduino nnoremap <buffer> <leader>au :ArduinoUpload<CR>
autocmd FileType arduino nnoremap <buffer> <leader>ad :ArduinoUploadAndSerial<CR>
autocmd FileType arduino nnoremap <buffer> <leader>ab :ArduinoChooseBoard<CR>
autocmd FileType arduino nnoremap <buffer> <leader>ap :ArduinoChooseProgrammer<CR>
autocmd FileType arduino nnoremap <buffer> <leader>r :ArduinoUpload<CR>
autocmd FileType arduino nnoremap <buffer> <leader>R :ArduinoUploadAndSerial<CR>

autocmd FileType haskell nnoremap <buffer> <leader>r :!ghc<Space>%<Space>&&<Space>./%:r<Enter>
autocmd FileType haskell nnoremap <buffer> <leader>R :vnew<Space>\|<Space>call<Space>StartGHCI()<Enter>

autocmd FileType c nnoremap <buffer> <leader>r :make! run<cr>

autocmd FileType python nnoremap <buffer> <Leader>r :!python<Space>%<Enter>
autocmd FileType tex nnoremap <buffer> <Leader>r :LLPStartPreview<Enter>
autocmd FileType sh nnoremap <buffer> <Leader>r :!./%<Enter>
autocmd FileType nroff nnoremap <buffer> <Leader>r :!/home/lukas/Documents/groff/compile.sh<space>%<Enter>
autocmd FileType vimwiki nnoremap <buffer> <Leader>r :VimwikiAll2HTML<Enter>
autocmd FileType dot nnoremap <buffer> <Leader>r :GraphvizInteractive<CR>
let vim_markdown_preview_hotkey='<Leader>r'

" ### Go specific
" Auto-Import dependencies on saving
let g:go_fmt_command = "goimports"

" Variable type info
let g:go_auto_type_info = 1

" breakpoints
autocmd FileType go nnoremap <buffer> <F9> :DlvToggleBreakpoint<Enter>
autocmd FileType go nnoremap <buffer> <F10> :DlvToggleTracepoint<Enter>

" ############# C Lang ###############

let &path.="src/include,/usr/include/AL,"

" ############## Ale code linting ##################
"
" Error and warning signs.
let g:ale_sign_error = '⤫'
let g:ale_sign_warning = '⚠'

" Enable integration with airline.
let g:airline#extensions#ale#enabled = 1

highlight clear ALEErrorSign
highlight clear ALEWarningSign

let g:ale_set_highlights = 0

" ### Markdown specific
let vim_markdown_preview_toggle=2
let vim_markdown_preview_browser='Firefox'
let vim_markdown_preview_github=1

" ################## Terminal ######################

" starts a terminal inside neovim
function StartTerminal()
	call OpenAutoCloseTerminal('bash')
	normal i
endfunction

" Start interactive haskell compiler inside neovim
function StartGHCI()
	call OpenAutoCloseTerminal('cd $(dirname %) && ghci')
	normal i
endfunction

" starts a file manager inside a terminal inside vim
function StartFileManager()
	call OpenAutoCloseTerminal('ranger')
	normal i
endfunction

function OpenAutoCloseTerminal(cmd)
	call termopen(a:cmd, {'on_exit': 'ExitTerminal'})
endfunction

function ExitTerminal(job_id, code, event)
	if a:code == 0
		close
	endif
endfunction

" Open terminal (shell) in vsplit window and start insert mode
nnoremap <Leader>s :vnew<Space>\|<Space>call<Space>StartTerminal()<Enter>
nnoremap <Leader>S :vnew<Space>\|<Space>call<Space>StartFileManager()<Enter>
autocmd BufEnter * if &buftype == 'terminal' | startinsert | endif
autocmd BufLeave * if &buftype == 'terminal' | stopinsert | endif

" Faster switching of split window focus
tnoremap <C-h> <C-\><C-N><C-w>h
tnoremap <C-j> <C-\><C-N><C-w>j
tnoremap <C-k> <C-\><C-N><C-w>k
tnoremap <C-l> <C-\><C-N><C-w>l
tnoremap <C-Left> <C-\><C-N><C-w>h
tnoremap <C-Down> <C-\><C-N><C-w>j
tnoremap <C-Up> <C-\><C-N><C-w>k
tnoremap <C-Right> <C-\><C-N><C-w>l

" ################# Startify ######################

nnoremap <Leader><C-s> :Startify<Enter>

let g:startify_enable_special      = 0
let g:startify_files_number        = 8
let g:startify_relative_path       = 1
let g:startify_change_to_dir       = 1
let g:startify_update_oldfiles     = 1
let g:startify_session_autoload    = 1
let g:startify_session_persistence = 1

let g:startify_skiplist = [
        \ 'COMMIT_EDITMSG',
        \ 'bundle/.*/doc',
        \ '/data/repo/neovim/runtime/doc',
        \ '/Users/mhi/local/vim/share/vim/vim74/doc',
        \ ]

let g:startify_bookmarks = [
		\ '~/workspace/',
        \ '~/go/src/github.com/phoenixdevelops/',
		\ '~/Documents/school/',
		\ '~/.config/',
		\ '~'
        \ ]

let g:startify_custom_header =
        \ startify#fortune#cowsay('', '═','║','╔','╗','╝','╚')


let g:startify_commands = [
    \ {'t': ['Start Terminal', 'new | call StartTerminal()']},
	\ {'f': ['Start Filemanager', 'new | call StartFileManager()']}
    \ ]

hi StartifyBracket ctermfg=240
hi StartifyFile    ctermfg=147
hi StartifyFooter  ctermfg=240
hi StartifyHeader  ctermfg=114
hi StartifyNumber  ctermfg=215
hi StartifyPath    ctermfg=245
hi StartifySlash   ctermfg=240
hi StartifySpecial ctermfg=240

" Exit session to startify screen
nnoremap <silent> <Leader>q :SClose<Enter>

" ################ Vim-Orgmode ####################

let g:org_aggressive_conceal = 1
let g:org_todo_keywords = ['TODO', 'NEXT', 'DONE']
let g:org_indent = 1

let g:org_todo_keyword_faces = [
	\ ['TODO', [':foreground red', ':background none', ':weight bold']],
	\ ['NEXT', [':foreground yellow', ':background none', 'weight bold']],
	\ ['DONE', [':foreground green', ':background none', 'wheight bold']]]

" ############## vim-arduino ######################

let g:arduino_dir = '/home/lukas/.arduino15'
let g:arduino_serial_baud = '19200'
let g:arduino_board = 'arduino:avr:nano'

" my_file.ino [arduino:avr:uno] [arduino:usbtinyisp] (/dev/ttyACM0:9600)
function! MyStatusLine()
  let port = arduino#GetPort()
  let line = '%f [' . g:arduino_board . '] [' . g:arduino_programmer . ']'
  if !empty(port)
    let line = line . ' (' . port . ':' . g:arduino_serial_baud . ')'
  endif
  return line
endfunction

autocmd FileType arduino let g:airline_section_x='%{MyStatusLine()}'

" ############### vimwiki ########################
let wiki = {}
let wiki.path = '~/vimwiki/'
let wiki.nested_syntaxes = {'c': 'c', 'go': 'go', 'xml': 'xml'}

let g:vimwiki_list = [{
  \ 'path': '$HOME/vimwiki',
  \ 'template_path': '$HOME/Templates/vimwiki',
  \ 'template_default': 'default',
  \ 'template_ext': '.html'}]

function NewDiaryEntry()
	let l:time = strftime('%Y-%m-%d_%H-%M')
	let l:path = g:vimwiki_list[0]['path']

	" construct file name
	let l:file = fnamemodify(path, ':r') . '/diary/' . time . '.wiki'

	" open the file
	exe 'e ' . l:file
endfunction

autocmd VimEnter * nnoremap <silent> <Leader>w<Leader>w :call<space>NewDiaryEntry()<CR>
