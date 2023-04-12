[34m[34m [i] [0m[0m [94m[94m# bash completion V2 for fleek                                -*- shell-script -*-[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m__fleek_debug()[0m[0m
[34m[34m     [0m[0m [94m[94m{[0m[0m
[34m[34m     [0m[0m [94m[94m    if [[ -n ${BASH_COMP_DEBUG_FILE-} ]]; then[0m[0m
[34m[34m     [0m[0m [94m[94m        echo "$*" >> "${BASH_COMP_DEBUG_FILE}"[0m[0m
[34m[34m     [0m[0m [94m[94m    fi[0m[0m
[34m[34m     [0m[0m [94m[94m}[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m# Macs have bash3 for which the bash-completion package doesn't include[0m[0m
[34m[34m     [0m[0m [94m[94m# _init_completion. This is a minimal version of that function.[0m[0m
[34m[34m     [0m[0m [94m[94m__fleek_init_completion()[0m[0m
[34m[34m     [0m[0m [94m[94m{[0m[0m
[34m[34m     [0m[0m [94m[94m    COMPREPLY=()[0m[0m
[34m[34m     [0m[0m [94m[94m    _get_comp_words_by_ref "$@" cur prev words cword[0m[0m
[34m[34m     [0m[0m [94m[94m}[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m# This function calls the fleek program to obtain the completion[0m[0m
[34m[34m     [0m[0m [94m[94m# results and the directive.  It fills the 'out' and 'directive' vars.[0m[0m
[34m[34m     [0m[0m [94m[94m__fleek_get_completion_results() {[0m[0m
[34m[34m     [0m[0m [94m[94m    local requestComp lastParam lastChar args[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    # Prepare the command to request completions for the program.[0m[0m
[34m[34m     [0m[0m [94m[94m    # Calling ${words[0]} instead of directly fleek allows to handle aliases[0m[0m
[34m[34m     [0m[0m [94m[94m    args=("${words[@]:1}")[0m[0m
[34m[34m     [0m[0m [94m[94m    requestComp="${words[0]} __complete ${args[*]}"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    lastParam=${words[$((${#words[@]}-1))]}[0m[0m
[34m[34m     [0m[0m [94m[94m    lastChar=${lastParam:$((${#lastParam}-1)):1}[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "lastParam ${lastParam}, lastChar ${lastChar}"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    if [[ -z ${cur} && ${lastChar} != = ]]; then[0m[0m
[34m[34m     [0m[0m [94m[94m        # If the last parameter is complete (there is a space following it)[0m[0m
[34m[34m     [0m[0m [94m[94m        # We add an extra empty parameter so we can indicate this to the go method.[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_debug "Adding extra empty parameter"[0m[0m
[34m[34m     [0m[0m [94m[94m        requestComp="${requestComp} ''"[0m[0m
[34m[34m     [0m[0m [94m[94m    fi[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    # When completing a flag with an = (e.g., fleek -n=<TAB>)[0m[0m
[34m[34m     [0m[0m [94m[94m    # bash focuses on the part after the =, so we need to remove[0m[0m
[34m[34m     [0m[0m [94m[94m    # the flag part from $cur[0m[0m
[34m[34m     [0m[0m [94m[94m    if [[ ${cur} == -*=* ]]; then[0m[0m
[34m[34m     [0m[0m [94m[94m        cur="${cur#*=}"[0m[0m
[34m[34m     [0m[0m [94m[94m    fi[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "Calling ${requestComp}"[0m[0m
[34m[34m     [0m[0m [94m[94m    # Use eval to handle any environment variables and such[0m[0m
[34m[34m     [0m[0m [94m[94m    out=$(eval "${requestComp}" 2>/dev/null)[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    # Extract the directive integer at the very end of the output following a colon (:)[0m[0m
[34m[34m     [0m[0m [94m[94m    directive=${out##*:}[0m[0m
[34m[34m     [0m[0m [94m[94m    # Remove the directive[0m[0m
[34m[34m     [0m[0m [94m[94m    out=${out%:*}[0m[0m
[34m[34m     [0m[0m [94m[94m    if [[ ${directive} == "${out}" ]]; then[0m[0m
[34m[34m     [0m[0m [94m[94m        # There is not directive specified[0m[0m
[34m[34m     [0m[0m [94m[94m        directive=0[0m[0m
[34m[34m     [0m[0m [94m[94m    fi[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "The completion directive is: ${directive}"[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "The completions are: ${out}"[0m[0m
[34m[34m     [0m[0m [94m[94m}[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m__fleek_process_completion_results() {[0m[0m
[34m[34m     [0m[0m [94m[94m    local shellCompDirectiveError=1[0m[0m
[34m[34m     [0m[0m [94m[94m    local shellCompDirectiveNoSpace=2[0m[0m
[34m[34m     [0m[0m [94m[94m    local shellCompDirectiveNoFileComp=4[0m[0m
[34m[34m     [0m[0m [94m[94m    local shellCompDirectiveFilterFileExt=8[0m[0m
[34m[34m     [0m[0m [94m[94m    local shellCompDirectiveFilterDirs=16[0m[0m
[34m[34m     [0m[0m [94m[94m    local shellCompDirectiveKeepOrder=32[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    if (((directive & shellCompDirectiveError) != 0)); then[0m[0m
[34m[34m     [0m[0m [94m[94m        # Error code.  No completion.[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_debug "Received error from custom completion go code"[0m[0m
[34m[34m     [0m[0m [94m[94m        return[0m[0m
[34m[34m     [0m[0m [94m[94m    else[0m[0m
[34m[34m     [0m[0m [94m[94m        if (((directive & shellCompDirectiveNoSpace) != 0)); then[0m[0m
[34m[34m     [0m[0m [94m[94m            if [[ $(type -t compopt) == builtin ]]; then[0m[0m
[34m[34m     [0m[0m [94m[94m                __fleek_debug "Activating no space"[0m[0m
[34m[34m     [0m[0m [94m[94m                compopt -o nospace[0m[0m
[34m[34m     [0m[0m [94m[94m            else[0m[0m
[34m[34m     [0m[0m [94m[94m                __fleek_debug "No space directive not supported in this version of bash"[0m[0m
[34m[34m     [0m[0m [94m[94m            fi[0m[0m
[34m[34m     [0m[0m [94m[94m        fi[0m[0m
[34m[34m     [0m[0m [94m[94m        if (((directive & shellCompDirectiveKeepOrder) != 0)); then[0m[0m
[34m[34m     [0m[0m [94m[94m            if [[ $(type -t compopt) == builtin ]]; then[0m[0m
[34m[34m     [0m[0m [94m[94m                # no sort isn't supported for bash less than < 4.4[0m[0m
[34m[34m     [0m[0m [94m[94m                if [[ ${BASH_VERSINFO[0]} -lt 4 || ( ${BASH_VERSINFO[0]} -eq 4 && ${BASH_VERSINFO[1]} -lt 4 ) ]]; then[0m[0m
[34m[34m     [0m[0m [94m[94m                    __fleek_debug "No sort directive not supported in this version of bash"[0m[0m
[34m[34m     [0m[0m [94m[94m                else[0m[0m
[34m[34m     [0m[0m [94m[94m                    __fleek_debug "Activating keep order"[0m[0m
[34m[34m     [0m[0m [94m[94m                    compopt -o nosort[0m[0m
[34m[34m     [0m[0m [94m[94m                fi[0m[0m
[34m[34m     [0m[0m [94m[94m            else[0m[0m
[34m[34m     [0m[0m [94m[94m                __fleek_debug "No sort directive not supported in this version of bash"[0m[0m
[34m[34m     [0m[0m [94m[94m            fi[0m[0m
[34m[34m     [0m[0m [94m[94m        fi[0m[0m
[34m[34m     [0m[0m [94m[94m        if (((directive & shellCompDirectiveNoFileComp) != 0)); then[0m[0m
[34m[34m     [0m[0m [94m[94m            if [[ $(type -t compopt) == builtin ]]; then[0m[0m
[34m[34m     [0m[0m [94m[94m                __fleek_debug "Activating no file completion"[0m[0m
[34m[34m     [0m[0m [94m[94m                compopt +o default[0m[0m
[34m[34m     [0m[0m [94m[94m            else[0m[0m
[34m[34m     [0m[0m [94m[94m                __fleek_debug "No file completion directive not supported in this version of bash"[0m[0m
[34m[34m     [0m[0m [94m[94m            fi[0m[0m
[34m[34m     [0m[0m [94m[94m        fi[0m[0m
[34m[34m     [0m[0m [94m[94m    fi[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    # Separate activeHelp from normal completions[0m[0m
[34m[34m     [0m[0m [94m[94m    local completions=()[0m[0m
[34m[34m     [0m[0m [94m[94m    local activeHelp=()[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_extract_activeHelp[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    if (((directive & shellCompDirectiveFilterFileExt) != 0)); then[0m[0m
[34m[34m     [0m[0m [94m[94m        # File extension filtering[0m[0m
[34m[34m     [0m[0m [94m[94m        local fullFilter filter filteringCmd[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m        # Do not use quotes around the $completions variable or else newline[0m[0m
[34m[34m     [0m[0m [94m[94m        # characters will be kept.[0m[0m
[34m[34m     [0m[0m [94m[94m        for filter in ${completions[*]}; do[0m[0m
[34m[34m     [0m[0m [94m[94m            fullFilter+="$filter|"[0m[0m
[34m[34m     [0m[0m [94m[94m        done[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m        filteringCmd="_filedir $fullFilter"[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_debug "File filtering command: $filteringCmd"[0m[0m
[34m[34m     [0m[0m [94m[94m        $filteringCmd[0m[0m
[34m[34m     [0m[0m [94m[94m    elif (((directive & shellCompDirectiveFilterDirs) != 0)); then[0m[0m
[34m[34m     [0m[0m [94m[94m        # File completion for directories only[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m        local subdir[0m[0m
[34m[34m     [0m[0m [94m[94m        subdir=${completions[0]}[0m[0m
[34m[34m     [0m[0m [94m[94m        if [[ -n $subdir ]]; then[0m[0m
[34m[34m     [0m[0m [94m[94m            __fleek_debug "Listing directories in $subdir"[0m[0m
[34m[34m     [0m[0m [94m[94m            pushd "$subdir" >/dev/null 2>&1 && _filedir -d && popd >/dev/null 2>&1 || return[0m[0m
[34m[34m     [0m[0m [94m[94m        else[0m[0m
[34m[34m     [0m[0m [94m[94m            __fleek_debug "Listing directories in ."[0m[0m
[34m[34m     [0m[0m [94m[94m            _filedir -d[0m[0m
[34m[34m     [0m[0m [94m[94m        fi[0m[0m
[34m[34m     [0m[0m [94m[94m    else[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_handle_completion_types[0m[0m
[34m[34m     [0m[0m [94m[94m    fi[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_handle_special_char "$cur" :[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_handle_special_char "$cur" =[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    # Print the activeHelp statements before we finish[0m[0m
[34m[34m     [0m[0m [94m[94m    if ((${#activeHelp[*]} != 0)); then[0m[0m
[34m[34m     [0m[0m [94m[94m        printf "\n";[0m[0m
[34m[34m     [0m[0m [94m[94m        printf "%s\n" "${activeHelp[@]}"[0m[0m
[34m[34m     [0m[0m [94m[94m        printf "\n"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m        # The prompt format is only available from bash 4.4.[0m[0m
[34m[34m     [0m[0m [94m[94m        # We test if it is available before using it.[0m[0m
[34m[34m     [0m[0m [94m[94m        if (x=${PS1@P}) 2> /dev/null; then[0m[0m
[34m[34m     [0m[0m [94m[94m            printf "%s" "${PS1@P}${COMP_LINE[@]}"[0m[0m
[34m[34m     [0m[0m [94m[94m        else[0m[0m
[34m[34m     [0m[0m [94m[94m            # Can't print the prompt.  Just print the[0m[0m
[34m[34m     [0m[0m [94m[94m            # text the user had typed, it is workable enough.[0m[0m
[34m[34m     [0m[0m [94m[94m            printf "%s" "${COMP_LINE[@]}"[0m[0m
[34m[34m     [0m[0m [94m[94m        fi[0m[0m
[34m[34m     [0m[0m [94m[94m    fi[0m[0m
[34m[34m     [0m[0m [94m[94m}[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m# Separate activeHelp lines from real completions.[0m[0m
[34m[34m     [0m[0m [94m[94m# Fills the $activeHelp and $completions arrays.[0m[0m
[34m[34m     [0m[0m [94m[94m__fleek_extract_activeHelp() {[0m[0m
[34m[34m     [0m[0m [94m[94m    local activeHelpMarker="_activeHelp_ "[0m[0m
[34m[34m     [0m[0m [94m[94m    local endIndex=${#activeHelpMarker}[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    while IFS='' read -r comp; do[0m[0m
[34m[34m     [0m[0m [94m[94m        if [[ ${comp:0:endIndex} == $activeHelpMarker ]]; then[0m[0m
[34m[34m     [0m[0m [94m[94m            comp=${comp:endIndex}[0m[0m
[34m[34m     [0m[0m [94m[94m            __fleek_debug "ActiveHelp found: $comp"[0m[0m
[34m[34m     [0m[0m [94m[94m            if [[ -n $comp ]]; then[0m[0m
[34m[34m     [0m[0m [94m[94m                activeHelp+=("$comp")[0m[0m
[34m[34m     [0m[0m [94m[94m            fi[0m[0m
[34m[34m     [0m[0m [94m[94m        else[0m[0m
[34m[34m     [0m[0m [94m[94m            # Not an activeHelp line but a normal completion[0m[0m
[34m[34m     [0m[0m [94m[94m            completions+=("$comp")[0m[0m
[34m[34m     [0m[0m [94m[94m        fi[0m[0m
[34m[34m     [0m[0m [94m[94m    done <<<"${out}"[0m[0m
[34m[34m     [0m[0m [94m[94m}[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m__fleek_handle_completion_types() {[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "__fleek_handle_completion_types: COMP_TYPE is $COMP_TYPE"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    case $COMP_TYPE in[0m[0m
[34m[34m     [0m[0m [94m[94m    37|42)[0m[0m
[34m[34m     [0m[0m [94m[94m        # Type: menu-complete/menu-complete-backward and insert-completions[0m[0m
[34m[34m     [0m[0m [94m[94m        # If the user requested inserting one completion at a time, or all[0m[0m
[34m[34m     [0m[0m [94m[94m        # completions at once on the command-line we must remove the descriptions.[0m[0m
[34m[34m     [0m[0m [94m[94m        # https://github.com/spf13/cobra/issues/1508[0m[0m
[34m[34m     [0m[0m [94m[94m        local tab=$'\t' comp[0m[0m
[34m[34m     [0m[0m [94m[94m        while IFS='' read -r comp; do[0m[0m
[34m[34m     [0m[0m [94m[94m            [[ -z $comp ]] && continue[0m[0m
[34m[34m     [0m[0m [94m[94m            # Strip any description[0m[0m
[34m[34m     [0m[0m [94m[94m            comp=${comp%%$tab*}[0m[0m
[34m[34m     [0m[0m [94m[94m            # Only consider the completions that match[0m[0m
[34m[34m     [0m[0m [94m[94m            if [[ $comp == "$cur"* ]]; then[0m[0m
[34m[34m     [0m[0m [94m[94m                COMPREPLY+=("$comp")[0m[0m
[34m[34m     [0m[0m [94m[94m            fi[0m[0m
[34m[34m     [0m[0m [94m[94m        done < <(printf "%s\n" "${completions[@]}")[0m[0m
[34m[34m     [0m[0m [94m[94m        ;;[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    *)[0m[0m
[34m[34m     [0m[0m [94m[94m        # Type: complete (normal completion)[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_handle_standard_completion_case[0m[0m
[34m[34m     [0m[0m [94m[94m        ;;[0m[0m
[34m[34m     [0m[0m [94m[94m    esac[0m[0m
[34m[34m     [0m[0m [94m[94m}[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m__fleek_handle_standard_completion_case() {[0m[0m
[34m[34m     [0m[0m [94m[94m    local tab=$'\t' comp[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    # Short circuit to optimize if we don't have descriptions[0m[0m
[34m[34m     [0m[0m [94m[94m    if [[ "${completions[*]}" != *$tab* ]]; then[0m[0m
[34m[34m     [0m[0m [94m[94m        IFS=$'\n' read -ra COMPREPLY -d '' < <(compgen -W "${completions[*]}" -- "$cur")[0m[0m
[34m[34m     [0m[0m [94m[94m        return 0[0m[0m
[34m[34m     [0m[0m [94m[94m    fi[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    local longest=0[0m[0m
[34m[34m     [0m[0m [94m[94m    local compline[0m[0m
[34m[34m     [0m[0m [94m[94m    # Look for the longest completion so that we can format things nicely[0m[0m
[34m[34m     [0m[0m [94m[94m    while IFS='' read -r compline; do[0m[0m
[34m[34m     [0m[0m [94m[94m        [[ -z $compline ]] && continue[0m[0m
[34m[34m     [0m[0m [94m[94m        # Strip any description before checking the length[0m[0m
[34m[34m     [0m[0m [94m[94m        comp=${compline%%$tab*}[0m[0m
[34m[34m     [0m[0m [94m[94m        # Only consider the completions that match[0m[0m
[34m[34m     [0m[0m [94m[94m        [[ $comp == "$cur"* ]] || continue[0m[0m
[34m[34m     [0m[0m [94m[94m        COMPREPLY+=("$compline")[0m[0m
[34m[34m     [0m[0m [94m[94m        if ((${#comp}>longest)); then[0m[0m
[34m[34m     [0m[0m [94m[94m            longest=${#comp}[0m[0m
[34m[34m     [0m[0m [94m[94m        fi[0m[0m
[34m[34m     [0m[0m [94m[94m    done < <(printf "%s\n" "${completions[@]}")[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    # If there is a single completion left, remove the description text[0m[0m
[34m[34m     [0m[0m [94m[94m    if ((${#COMPREPLY[*]} == 1)); then[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_debug "COMPREPLY[0]: ${COMPREPLY[0]}"[0m[0m
[34m[34m     [0m[0m [94m[94m        comp="${COMPREPLY[0]%%$tab*}"[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_debug "Removed description from single completion, which is now: ${comp}"[0m[0m
[34m[34m     [0m[0m [94m[94m        COMPREPLY[0]=$comp[0m[0m
[34m[34m     [0m[0m [94m[94m    else # Format the descriptions[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_format_comp_descriptions $longest[0m[0m
[34m[34m     [0m[0m [94m[94m    fi[0m[0m
[34m[34m     [0m[0m [94m[94m}[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m__fleek_handle_special_char()[0m[0m
[34m[34m     [0m[0m [94m[94m{[0m[0m
[34m[34m     [0m[0m [94m[94m    local comp="$1"[0m[0m
[34m[34m     [0m[0m [94m[94m    local char=$2[0m[0m
[34m[34m     [0m[0m [94m[94m    if [[ "$comp" == *${char}* && "$COMP_WORDBREAKS" == *${char}* ]]; then[0m[0m
[34m[34m     [0m[0m [94m[94m        local word=${comp%"${comp##*${char}}"}[0m[0m
[34m[34m     [0m[0m [94m[94m        local idx=${#COMPREPLY[*]}[0m[0m
[34m[34m     [0m[0m [94m[94m        while ((--idx >= 0)); do[0m[0m
[34m[34m     [0m[0m [94m[94m            COMPREPLY[idx]=${COMPREPLY[idx]#"$word"}[0m[0m
[34m[34m     [0m[0m [94m[94m        done[0m[0m
[34m[34m     [0m[0m [94m[94m    fi[0m[0m
[34m[34m     [0m[0m [94m[94m}[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m__fleek_format_comp_descriptions()[0m[0m
[34m[34m     [0m[0m [94m[94m{[0m[0m
[34m[34m     [0m[0m [94m[94m    local tab=$'\t'[0m[0m
[34m[34m     [0m[0m [94m[94m    local comp desc maxdesclength[0m[0m
[34m[34m     [0m[0m [94m[94m    local longest=$1[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    local i ci[0m[0m
[34m[34m     [0m[0m [94m[94m    for ci in ${!COMPREPLY[*]}; do[0m[0m
[34m[34m     [0m[0m [94m[94m        comp=${COMPREPLY[ci]}[0m[0m
[34m[34m     [0m[0m [94m[94m        # Properly format the description string which follows a tab character if there is one[0m[0m
[34m[34m     [0m[0m [94m[94m        if [[ "$comp" == *$tab* ]]; then[0m[0m
[34m[34m     [0m[0m [94m[94m            __fleek_debug "Original comp: $comp"[0m[0m
[34m[34m     [0m[0m [94m[94m            desc=${comp#*$tab}[0m[0m
[34m[34m     [0m[0m [94m[94m            comp=${comp%%$tab*}[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m            # $COLUMNS stores the current shell width.[0m[0m
[34m[34m     [0m[0m [94m[94m            # Remove an extra 4 because we add 2 spaces and 2 parentheses.[0m[0m
[34m[34m     [0m[0m [94m[94m            maxdesclength=$(( COLUMNS - longest - 4 ))[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m            # Make sure we can fit a description of at least 8 characters[0m[0m
[34m[34m     [0m[0m [94m[94m            # if we are to align the descriptions.[0m[0m
[34m[34m     [0m[0m [94m[94m            if ((maxdesclength > 8)); then[0m[0m
[34m[34m     [0m[0m [94m[94m                # Add the proper number of spaces to align the descriptions[0m[0m
[34m[34m     [0m[0m [94m[94m                for ((i = ${#comp} ; i < longest ; i++)); do[0m[0m
[34m[34m     [0m[0m [94m[94m                    comp+=" "[0m[0m
[34m[34m     [0m[0m [94m[94m                done[0m[0m
[34m[34m     [0m[0m [94m[94m            else[0m[0m
[34m[34m     [0m[0m [94m[94m                # Don't pad the descriptions so we can fit more text after the completion[0m[0m
[34m[34m     [0m[0m [94m[94m                maxdesclength=$(( COLUMNS - ${#comp} - 4 ))[0m[0m
[34m[34m     [0m[0m [94m[94m            fi[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m            # If there is enough space for any description text,[0m[0m
[34m[34m     [0m[0m [94m[94m            # truncate the descriptions that are too long for the shell width[0m[0m
[34m[34m     [0m[0m [94m[94m            if ((maxdesclength > 0)); then[0m[0m
[34m[34m     [0m[0m [94m[94m                if ((${#desc} > maxdesclength)); then[0m[0m
[34m[34m     [0m[0m [94m[94m                    desc=${desc:0:$(( maxdesclength - 1 ))}[0m[0m
[34m[34m     [0m[0m [94m[94m                    desc+="…"[0m[0m
[34m[34m     [0m[0m [94m[94m                fi[0m[0m
[34m[34m     [0m[0m [94m[94m                comp+="  ($desc)"[0m[0m
[34m[34m     [0m[0m [94m[94m            fi[0m[0m
[34m[34m     [0m[0m [94m[94m            COMPREPLY[ci]=$comp[0m[0m
[34m[34m     [0m[0m [94m[94m            __fleek_debug "Final comp: $comp"[0m[0m
[34m[34m     [0m[0m [94m[94m        fi[0m[0m
[34m[34m     [0m[0m [94m[94m    done[0m[0m
[34m[34m     [0m[0m [94m[94m}[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m__start_fleek()[0m[0m
[34m[34m     [0m[0m [94m[94m{[0m[0m
[34m[34m     [0m[0m [94m[94m    local cur prev words cword split[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    COMPREPLY=()[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    # Call _init_completion from the bash-completion package[0m[0m
[34m[34m     [0m[0m [94m[94m    # to prepare the arguments properly[0m[0m
[34m[34m     [0m[0m [94m[94m    if declare -F _init_completion >/dev/null 2>&1; then[0m[0m
[34m[34m     [0m[0m [94m[94m        _init_completion -n =: || return[0m[0m
[34m[34m     [0m[0m [94m[94m    else[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_init_completion -n =: || return[0m[0m
[34m[34m     [0m[0m [94m[94m    fi[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "========= starting completion logic =========="[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "cur is ${cur}, words[*] is ${words[*]}, #words[@] is ${#words[@]}, cword is $cword"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    # The user could have moved the cursor backwards on the command-line.[0m[0m
[34m[34m     [0m[0m [94m[94m    # We need to trigger completion from the $cword location, so we need[0m[0m
[34m[34m     [0m[0m [94m[94m    # to truncate the command-line ($words) up to the $cword location.[0m[0m
[34m[34m     [0m[0m [94m[94m    words=("${words[@]:0:$cword+1}")[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "Truncated words[*]: ${words[*]},"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    local out directive[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_get_completion_results[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_process_completion_results[0m[0m
[34m[34m     [0m[0m [94m[94m}[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94mif [[ $(type -t compopt) = "builtin" ]]; then[0m[0m
[34m[34m     [0m[0m [94m[94m    complete -o default -F __start_fleek fleek[0m[0m
[34m[34m     [0m[0m [94m[94melse[0m[0m
[34m[34m     [0m[0m [94m[94m    complete -o default -o nospace -F __start_fleek fleek[0m[0m
[34m[34m     [0m[0m [94m[94mfi[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m# ex: ts=4 sw=4 et filetype=sh[0m[0m
