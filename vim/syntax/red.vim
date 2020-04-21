set syntax=sh

syntax keyword httpMethods 
    \ GET
    \ POST
    \ PUT
    \ DELETE
    \ OPTIONS

syntax match redAttrs "\v\@.+$"


highlight default link redAttrs PreProc
highlight default link httpMethods Keyword
