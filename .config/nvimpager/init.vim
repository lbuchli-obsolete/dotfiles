call plug#begin('~/.config/nvim/bundle')
" List of all plugins to use
Plug 'scrooloose/syntastic' " Syntax highlighting
Plug 'sheerun/vim-polyglot' " Language support for some languages
Plug 'SirVer/ultisnips' " Snippets for autocompletion, ...
Plug 'honza/vim-snippets' " Actual snippets for ultisnips
call plug#end()

" Initialization

" Basic variables
filetype plugin indent on
syntax on
let mapleader = "\<space>"
set number
set incsearch
set nohlsearch
set smartcase
set tabstop=4
set softtabstop=0
set expandtab
set shiftwidth=4
set noswapfile
set nowrap

" Make split windows look more clean
set fillchars+=vert:\ 
hi VertSplit ctermfg=LightGray

"#######Preferences########
        
" Vim
" noremap Y 0y$
set hidden
set history=100
" Cancel searches with Escape
nnoremap <silent> <Esc> :nohlsearch<Bar>:echo<CR>
" Highlight trailing whitespaces
autocmd ColorScheme * highlight ExtraWhitespace ctermbg=grey guibg=grey
" get rid of ~
"highlight EndOfBuffer ctermfg=bg ctermbg=NONE

" Autocompletion colors
highlight Pmenu ctermbg=0 ctermfg=15
highlight PmenuSel ctermbg=7 ctermfg=15
highlight PmenuSbar ctermbg=8 ctermfg=15

hi MatchParen cterm=bold ctermbg=None ctermfg=Red

" ################# Running Programs ##################
autocmd FileType go nnoremap <buffer> <Leader>r :!go<Space>run<Space>%<CR> 
autocmd FileType python nnoremap <buffer> <Leader>r :!python<Space>%<CR>
autocmd FileType tex nnoremap <buffer> <Leader>r :LLPStartPreview<CR>
autocmd FileType sh nnoremap <buffer> <Leader>r :!./%<CR>
