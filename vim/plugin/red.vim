" Red utility

let g:red_query_out = '$HOME/.red/query.redout'

" Runs query
command! -range=% RedQuery :<line1>,<line2>!red -c run

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
