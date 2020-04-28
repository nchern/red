" Red utility

let g:red_query_out = '$HOME/.red/query.redout'

function! RedDoQuery() range
    " Do some more things
    let lines = getline(a:firstline, a:lastline)
    " let selectedText = join(lines, "\n")
    " echo selectedText

    let out = system("red -c run -d=false -o - -s ". expand("%:p"), l:lines)

    let buf_name = "__RED_Results__"

    let win_num=bufwinnr(l:buf_name)
    if l:win_num < 0
        " Open a new split and set it up.
        exe "rightbelow vsplit" l:buf_name
        setlocal filetype=redout
        setlocal buftype=nofile
    else
        exe l:win_num . "wincmd w"
    endif

    " clear old stuff
    normal! ggdG
    " insert results
    call append(0, split(l:out, '\v\n'))
endfunction

" Runs query
command! -range=% RedQuery <line1>,<line2>call RedDoQuery()
" Runs query - backwards compatibility
command! -range=% RedQueryOld :<line1>,<line2>!red -c run

command! RedShowResults :execute 'botright vsp ' . g:red_query_out

" Only do the rest when the FileType autocommand has not been triggered yet.
if did_filetype()
  finish
endif

let s:line1 = getline(1)

if s:line1 =~# "^#!red"
    set syntax=red
endif

" TODO: do we need toi couple JSON formatting here?
" FIXME: does not behave nice if json contains erros
" let g:red_json_formatter = 'jq'
" command! -range RedFmtJSON :execute '<line1>,<line2>!' . g:red_json_formatter . ' . '
