[34m[34m [i] [0m[0m [94m[94m#compdef fleek[0m[0m
[34m[34m     [0m[0m [94m[94mcompdef _fleek fleek[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m# zsh completion for fleek                                -*- shell-script -*-[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m__fleek_debug()[0m[0m
[34m[34m     [0m[0m [94m[94m{[0m[0m
[34m[34m     [0m[0m [94m[94m    local file="$BASH_COMP_DEBUG_FILE"[0m[0m
[34m[34m     [0m[0m [94m[94m    if [[ -n ${file} ]]; then[0m[0m
[34m[34m     [0m[0m [94m[94m        echo "$*" >> "${file}"[0m[0m
[34m[34m     [0m[0m [94m[94m    fi[0m[0m
[34m[34m     [0m[0m [94m[94m}[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m_fleek()[0m[0m
[34m[34m     [0m[0m [94m[94m{[0m[0m
[34m[34m     [0m[0m [94m[94m    local shellCompDirectiveError=1[0m[0m
[34m[34m     [0m[0m [94m[94m    local shellCompDirectiveNoSpace=2[0m[0m
[34m[34m     [0m[0m [94m[94m    local shellCompDirectiveNoFileComp=4[0m[0m
[34m[34m     [0m[0m [94m[94m    local shellCompDirectiveFilterFileExt=8[0m[0m
[34m[34m     [0m[0m [94m[94m    local shellCompDirectiveFilterDirs=16[0m[0m
[34m[34m     [0m[0m [94m[94m    local shellCompDirectiveKeepOrder=32[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    local lastParam lastChar flagPrefix requestComp out directive comp lastComp noSpace keepOrder[0m[0m
[34m[34m     [0m[0m [94m[94m    local -a completions[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "\n========= starting completion logic =========="[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "CURRENT: ${CURRENT}, words[*]: ${words[*]}"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    # The user could have moved the cursor backwards on the command-line.[0m[0m
[34m[34m     [0m[0m [94m[94m    # We need to trigger completion from the $CURRENT location, so we need[0m[0m
[34m[34m     [0m[0m [94m[94m    # to truncate the command-line ($words) up to the $CURRENT location.[0m[0m
[34m[34m     [0m[0m [94m[94m    # (We cannot use $CURSOR as its value does not work when a command is an alias.)[0m[0m
[34m[34m     [0m[0m [94m[94m    words=("${=words[1,CURRENT]}")[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "Truncated words[*]: ${words[*]},"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    lastParam=${words[-1]}[0m[0m
[34m[34m     [0m[0m [94m[94m    lastChar=${lastParam[-1]}[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "lastParam: ${lastParam}, lastChar: ${lastChar}"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    # For zsh, when completing a flag with an = (e.g., fleek -n=<TAB>)[0m[0m
[34m[34m     [0m[0m [94m[94m    # completions must be prefixed with the flag[0m[0m
[34m[34m     [0m[0m [94m[94m    setopt local_options BASH_REMATCH[0m[0m
[34m[34m     [0m[0m [94m[94m    if [[ "${lastParam}" =~ '-.*=' ]]; then[0m[0m
[34m[34m     [0m[0m [94m[94m        # We are dealing with a flag with an =[0m[0m
[34m[34m     [0m[0m [94m[94m        flagPrefix="-P ${BASH_REMATCH}"[0m[0m
[34m[34m     [0m[0m [94m[94m    fi[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    # Prepare the command to obtain completions[0m[0m
[34m[34m     [0m[0m [94m[94m    requestComp="${words[1]} __complete ${words[2,-1]}"[0m[0m
[34m[34m     [0m[0m [94m[94m    if [ "${lastChar}" = "" ]; then[0m[0m
[34m[34m     [0m[0m [94m[94m        # If the last parameter is complete (there is a space following it)[0m[0m
[34m[34m     [0m[0m [94m[94m        # We add an extra empty parameter so we can indicate this to the go completion code.[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_debug "Adding extra empty parameter"[0m[0m
[34m[34m     [0m[0m [94m[94m        requestComp="${requestComp} \"\""[0m[0m
[34m[34m     [0m[0m [94m[94m    fi[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "About to call: eval ${requestComp}"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    # Use eval to handle any environment variables and such[0m[0m
[34m[34m     [0m[0m [94m[94m    out=$(eval ${requestComp} 2>/dev/null)[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "completion output: ${out}"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    # Extract the directive integer following a : from the last line[0m[0m
[34m[34m     [0m[0m [94m[94m    local lastLine[0m[0m
[34m[34m     [0m[0m [94m[94m    while IFS='\n' read -r line; do[0m[0m
[34m[34m     [0m[0m [94m[94m        lastLine=${line}[0m[0m
[34m[34m     [0m[0m [94m[94m    done < <(printf "%s\n" "${out[@]}")[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "last line: ${lastLine}"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    if [ "${lastLine[1]}" = : ]; then[0m[0m
[34m[34m     [0m[0m [94m[94m        directive=${lastLine[2,-1]}[0m[0m
[34m[34m     [0m[0m [94m[94m        # Remove the directive including the : and the newline[0m[0m
[34m[34m     [0m[0m [94m[94m        local suffix[0m[0m
[34m[34m     [0m[0m [94m[94m        (( suffix=${#lastLine}+2))[0m[0m
[34m[34m     [0m[0m [94m[94m        out=${out[1,-$suffix]}[0m[0m
[34m[34m     [0m[0m [94m[94m    else[0m[0m
[34m[34m     [0m[0m [94m[94m        # There is no directive specified.  Leave $out as is.[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_debug "No directive found.  Setting do default"[0m[0m
[34m[34m     [0m[0m [94m[94m        directive=0[0m[0m
[34m[34m     [0m[0m [94m[94m    fi[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "directive: ${directive}"[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "completions: ${out}"[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "flagPrefix: ${flagPrefix}"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    if [ $((directive & shellCompDirectiveError)) -ne 0 ]; then[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_debug "Completion received error. Ignoring completions."[0m[0m
[34m[34m     [0m[0m [94m[94m        return[0m[0m
[34m[34m     [0m[0m [94m[94m    fi[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    local activeHelpMarker="_activeHelp_ "[0m[0m
[34m[34m     [0m[0m [94m[94m    local endIndex=${#activeHelpMarker}[0m[0m
[34m[34m     [0m[0m [94m[94m    local startIndex=$((${#activeHelpMarker}+1))[0m[0m
[34m[34m     [0m[0m [94m[94m    local hasActiveHelp=0[0m[0m
[34m[34m     [0m[0m [94m[94m    while IFS='\n' read -r comp; do[0m[0m
[34m[34m     [0m[0m [94m[94m        # Check if this is an activeHelp statement (i.e., prefixed with $activeHelpMarker)[0m[0m
[34m[34m     [0m[0m [94m[94m        if [ "${comp[1,$endIndex]}" = "$activeHelpMarker" ];then[0m[0m
[34m[34m     [0m[0m [94m[94m            __fleek_debug "ActiveHelp found: $comp"[0m[0m
[34m[34m     [0m[0m [94m[94m            comp="${comp[$startIndex,-1]}"[0m[0m
[34m[34m     [0m[0m [94m[94m            if [ -n "$comp" ]; then[0m[0m
[34m[34m     [0m[0m [94m[94m                compadd -x "${comp}"[0m[0m
[34m[34m     [0m[0m [94m[94m                __fleek_debug "ActiveHelp will need delimiter"[0m[0m
[34m[34m     [0m[0m [94m[94m                hasActiveHelp=1[0m[0m
[34m[34m     [0m[0m [94m[94m            fi[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m            continue[0m[0m
[34m[34m     [0m[0m [94m[94m        fi[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m        if [ -n "$comp" ]; then[0m[0m
[34m[34m     [0m[0m [94m[94m            # If requested, completions are returned with a description.[0m[0m
[34m[34m     [0m[0m [94m[94m            # The description is preceded by a TAB character.[0m[0m
[34m[34m     [0m[0m [94m[94m            # For zsh's _describe, we need to use a : instead of a TAB.[0m[0m
[34m[34m     [0m[0m [94m[94m            # We first need to escape any : as part of the completion itself.[0m[0m
[34m[34m     [0m[0m [94m[94m            comp=${comp//:/\\:}[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m            local tab="$(printf '\t')"[0m[0m
[34m[34m     [0m[0m [94m[94m            comp=${comp//$tab/:}[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m            __fleek_debug "Adding completion: ${comp}"[0m[0m
[34m[34m     [0m[0m [94m[94m            completions+=${comp}[0m[0m
[34m[34m     [0m[0m [94m[94m            lastComp=$comp[0m[0m
[34m[34m     [0m[0m [94m[94m        fi[0m[0m
[34m[34m     [0m[0m [94m[94m    done < <(printf "%s\n" "${out[@]}")[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    # Add a delimiter after the activeHelp statements, but only if:[0m[0m
[34m[34m     [0m[0m [94m[94m    # - there are completions following the activeHelp statements, or[0m[0m
[34m[34m     [0m[0m [94m[94m    # - file completion will be performed (so there will be choices after the activeHelp)[0m[0m
[34m[34m     [0m[0m [94m[94m    if [ $hasActiveHelp -eq 1 ]; then[0m[0m
[34m[34m     [0m[0m [94m[94m        if [ ${#completions} -ne 0 ] || [ $((directive & shellCompDirectiveNoFileComp)) -eq 0 ]; then[0m[0m
[34m[34m     [0m[0m [94m[94m            __fleek_debug "Adding activeHelp delimiter"[0m[0m
[34m[34m     [0m[0m [94m[94m            compadd -x "--"[0m[0m
[34m[34m     [0m[0m [94m[94m            hasActiveHelp=0[0m[0m
[34m[34m     [0m[0m [94m[94m        fi[0m[0m
[34m[34m     [0m[0m [94m[94m    fi[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    if [ $((directive & shellCompDirectiveNoSpace)) -ne 0 ]; then[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_debug "Activating nospace."[0m[0m
[34m[34m     [0m[0m [94m[94m        noSpace="-S ''"[0m[0m
[34m[34m     [0m[0m [94m[94m    fi[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    if [ $((directive & shellCompDirectiveKeepOrder)) -ne 0 ]; then[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_debug "Activating keep order."[0m[0m
[34m[34m     [0m[0m [94m[94m        keepOrder="-V"[0m[0m
[34m[34m     [0m[0m [94m[94m    fi[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    if [ $((directive & shellCompDirectiveFilterFileExt)) -ne 0 ]; then[0m[0m
[34m[34m     [0m[0m [94m[94m        # File extension filtering[0m[0m
[34m[34m     [0m[0m [94m[94m        local filteringCmd[0m[0m
[34m[34m     [0m[0m [94m[94m        filteringCmd='_files'[0m[0m
[34m[34m     [0m[0m [94m[94m        for filter in ${completions[@]}; do[0m[0m
[34m[34m     [0m[0m [94m[94m            if [ ${filter[1]} != '*' ]; then[0m[0m
[34m[34m     [0m[0m [94m[94m                # zsh requires a glob pattern to do file filtering[0m[0m
[34m[34m     [0m[0m [94m[94m                filter="\*.$filter"[0m[0m
[34m[34m     [0m[0m [94m[94m            fi[0m[0m
[34m[34m     [0m[0m [94m[94m            filteringCmd+=" -g $filter"[0m[0m
[34m[34m     [0m[0m [94m[94m        done[0m[0m
[34m[34m     [0m[0m [94m[94m        filteringCmd+=" ${flagPrefix}"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_debug "File filtering command: $filteringCmd"[0m[0m
[34m[34m     [0m[0m [94m[94m        _arguments '*:filename:'"$filteringCmd"[0m[0m
[34m[34m     [0m[0m [94m[94m    elif [ $((directive & shellCompDirectiveFilterDirs)) -ne 0 ]; then[0m[0m
[34m[34m     [0m[0m [94m[94m        # File completion for directories only[0m[0m
[34m[34m     [0m[0m [94m[94m        local subdir[0m[0m
[34m[34m     [0m[0m [94m[94m        subdir="${completions[1]}"[0m[0m
[34m[34m     [0m[0m [94m[94m        if [ -n "$subdir" ]; then[0m[0m
[34m[34m     [0m[0m [94m[94m            __fleek_debug "Listing directories in $subdir"[0m[0m
[34m[34m     [0m[0m [94m[94m            pushd "${subdir}" >/dev/null 2>&1[0m[0m
[34m[34m     [0m[0m [94m[94m        else[0m[0m
[34m[34m     [0m[0m [94m[94m            __fleek_debug "Listing directories in ."[0m[0m
[34m[34m     [0m[0m [94m[94m        fi[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m        local result[0m[0m
[34m[34m     [0m[0m [94m[94m        _arguments '*:dirname:_files -/'" ${flagPrefix}"[0m[0m
[34m[34m     [0m[0m [94m[94m        result=$?[0m[0m
[34m[34m     [0m[0m [94m[94m        if [ -n "$subdir" ]; then[0m[0m
[34m[34m     [0m[0m [94m[94m            popd >/dev/null 2>&1[0m[0m
[34m[34m     [0m[0m [94m[94m        fi[0m[0m
[34m[34m     [0m[0m [94m[94m        return $result[0m[0m
[34m[34m     [0m[0m [94m[94m    else[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_debug "Calling _describe"[0m[0m
[34m[34m     [0m[0m [94m[94m        if eval _describe $keepOrder "completions" completions $flagPrefix $noSpace; then[0m[0m
[34m[34m     [0m[0m [94m[94m            __fleek_debug "_describe found some completions"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m            # Return the success of having called _describe[0m[0m
[34m[34m     [0m[0m [94m[94m            return 0[0m[0m
[34m[34m     [0m[0m [94m[94m        else[0m[0m
[34m[34m     [0m[0m [94m[94m            __fleek_debug "_describe did not find completions."[0m[0m
[34m[34m     [0m[0m [94m[94m            __fleek_debug "Checking if we should do file completion."[0m[0m
[34m[34m     [0m[0m [94m[94m            if [ $((directive & shellCompDirectiveNoFileComp)) -ne 0 ]; then[0m[0m
[34m[34m     [0m[0m [94m[94m                __fleek_debug "deactivating file completion"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m                # We must return an error code here to let zsh know that there were no[0m[0m
[34m[34m     [0m[0m [94m[94m                # completions found by _describe; this is what will trigger other[0m[0m
[34m[34m     [0m[0m [94m[94m                # matching algorithms to attempt to find completions.[0m[0m
[34m[34m     [0m[0m [94m[94m                # For example zsh can match letters in the middle of words.[0m[0m
[34m[34m     [0m[0m [94m[94m                return 1[0m[0m
[34m[34m     [0m[0m [94m[94m            else[0m[0m
[34m[34m     [0m[0m [94m[94m                # Perform file completion[0m[0m
[34m[34m     [0m[0m [94m[94m                __fleek_debug "Activating file completion"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m                # We must return the result of this command, so it must be the[0m[0m
[34m[34m     [0m[0m [94m[94m                # last command, or else we must store its result to return it.[0m[0m
[34m[34m     [0m[0m [94m[94m                _arguments '*:filename:_files'" ${flagPrefix}"[0m[0m
[34m[34m     [0m[0m [94m[94m            fi[0m[0m
[34m[34m     [0m[0m [94m[94m        fi[0m[0m
[34m[34m     [0m[0m [94m[94m    fi[0m[0m
[34m[34m     [0m[0m [94m[94m}[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m# don't run the completion function when being source-ed or eval-ed[0m[0m
[34m[34m     [0m[0m [94m[94mif [ "$funcstack[1]" = "_fleek" ]; then[0m[0m
[34m[34m     [0m[0m [94m[94m    _fleek[0m[0m
[34m[34m     [0m[0m [94m[94mfi[0m[0m
