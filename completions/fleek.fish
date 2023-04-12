[34m[34m [i] [0m[0m [94m[94m# fish completion for fleek                                -*- shell-script -*-[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94mfunction __fleek_debug[0m[0m
[34m[34m     [0m[0m [94m[94m    set -l file "$BASH_COMP_DEBUG_FILE"[0m[0m
[34m[34m     [0m[0m [94m[94m    if test -n "$file"[0m[0m
[34m[34m     [0m[0m [94m[94m        echo "$argv" >> $file[0m[0m
[34m[34m     [0m[0m [94m[94m    end[0m[0m
[34m[34m     [0m[0m [94m[94mend[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94mfunction __fleek_perform_completion[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "Starting __fleek_perform_completion"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    # Extract all args except the last one[0m[0m
[34m[34m     [0m[0m [94m[94m    set -l args (commandline -opc)[0m[0m
[34m[34m     [0m[0m [94m[94m    # Extract the last arg and escape it in case it is a space[0m[0m
[34m[34m     [0m[0m [94m[94m    set -l lastArg (string escape -- (commandline -ct))[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "args: $args"[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "last arg: $lastArg"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    # Disable ActiveHelp which is not supported for fish shell[0m[0m
[34m[34m     [0m[0m [94m[94m    set -l requestComp "FLEEK_ACTIVE_HELP=0 $args[1] __complete $args[2..-1] $lastArg"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "Calling $requestComp"[0m[0m
[34m[34m     [0m[0m [94m[94m    set -l results (eval $requestComp 2> /dev/null)[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    # Some programs may output extra empty lines after the directive.[0m[0m
[34m[34m     [0m[0m [94m[94m    # Let's ignore them or else it will break completion.[0m[0m
[34m[34m     [0m[0m [94m[94m    # Ref: https://github.com/spf13/cobra/issues/1279[0m[0m
[34m[34m     [0m[0m [94m[94m    for line in $results[-1..1][0m[0m
[34m[34m     [0m[0m [94m[94m        if test (string trim -- $line) = ""[0m[0m
[34m[34m     [0m[0m [94m[94m            # Found an empty line, remove it[0m[0m
[34m[34m     [0m[0m [94m[94m            set results $results[1..-2][0m[0m
[34m[34m     [0m[0m [94m[94m        else[0m[0m
[34m[34m     [0m[0m [94m[94m            # Found non-empty line, we have our proper output[0m[0m
[34m[34m     [0m[0m [94m[94m            break[0m[0m
[34m[34m     [0m[0m [94m[94m        end[0m[0m
[34m[34m     [0m[0m [94m[94m    end[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    set -l comps $results[1..-2][0m[0m
[34m[34m     [0m[0m [94m[94m    set -l directiveLine $results[-1][0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    # For Fish, when completing a flag with an = (e.g., <program> -n=<TAB>)[0m[0m
[34m[34m     [0m[0m [94m[94m    # completions must be prefixed with the flag[0m[0m
[34m[34m     [0m[0m [94m[94m    set -l flagPrefix (string match -r -- '-.*=' "$lastArg")[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "Comps: $comps"[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "DirectiveLine: $directiveLine"[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "flagPrefix: $flagPrefix"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    for comp in $comps[0m[0m
[34m[34m     [0m[0m [94m[94m        printf "%s%s\n" "$flagPrefix" "$comp"[0m[0m
[34m[34m     [0m[0m [94m[94m    end[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    printf "%s\n" "$directiveLine"[0m[0m
[34m[34m     [0m[0m [94m[94mend[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m# this function limits calls to __fleek_perform_completion, by caching the result behind $__fleek_perform_completion_once_result[0m[0m
[34m[34m     [0m[0m [94m[94mfunction __fleek_perform_completion_once[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "Starting __fleek_perform_completion_once"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    if test -n "$__fleek_perform_completion_once_result"[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_debug "Seems like a valid result already exists, skipping __fleek_perform_completion"[0m[0m
[34m[34m     [0m[0m [94m[94m        return 0[0m[0m
[34m[34m     [0m[0m [94m[94m    end[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    set --global __fleek_perform_completion_once_result (__fleek_perform_completion)[0m[0m
[34m[34m     [0m[0m [94m[94m    if test -z "$__fleek_perform_completion_once_result"[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_debug "No completions, probably due to a failure"[0m[0m
[34m[34m     [0m[0m [94m[94m        return 1[0m[0m
[34m[34m     [0m[0m [94m[94m    end[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "Performed completions and set __fleek_perform_completion_once_result"[0m[0m
[34m[34m     [0m[0m [94m[94m    return 0[0m[0m
[34m[34m     [0m[0m [94m[94mend[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m# this function is used to clear the $__fleek_perform_completion_once_result variable after completions are run[0m[0m
[34m[34m     [0m[0m [94m[94mfunction __fleek_clear_perform_completion_once_result[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug ""[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "========= clearing previously set __fleek_perform_completion_once_result variable =========="[0m[0m
[34m[34m     [0m[0m [94m[94m    set --erase __fleek_perform_completion_once_result[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "Succesfully erased the variable __fleek_perform_completion_once_result"[0m[0m
[34m[34m     [0m[0m [94m[94mend[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94mfunction __fleek_requires_order_preservation[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug ""[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "========= checking if order preservation is required =========="[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_perform_completion_once[0m[0m
[34m[34m     [0m[0m [94m[94m    if test -z "$__fleek_perform_completion_once_result"[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_debug "Error determining if order preservation is required"[0m[0m
[34m[34m     [0m[0m [94m[94m        return 1[0m[0m
[34m[34m     [0m[0m [94m[94m    end[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    set -l directive (string sub --start 2 $__fleek_perform_completion_once_result[-1])[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "Directive is: $directive"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    set -l shellCompDirectiveKeepOrder 32[0m[0m
[34m[34m     [0m[0m [94m[94m    set -l keeporder (math (math --scale 0 $directive / $shellCompDirectiveKeepOrder) % 2)[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "Keeporder is: $keeporder"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    if test $keeporder -ne 0[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_debug "This does require order preservation"[0m[0m
[34m[34m     [0m[0m [94m[94m        return 0[0m[0m
[34m[34m     [0m[0m [94m[94m    end[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "This doesn't require order preservation"[0m[0m
[34m[34m     [0m[0m [94m[94m    return 1[0m[0m
[34m[34m     [0m[0m [94m[94mend[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m# This function does two things:[0m[0m
[34m[34m     [0m[0m [94m[94m# - Obtain the completions and store them in the global __fleek_comp_results[0m[0m
[34m[34m     [0m[0m [94m[94m# - Return false if file completion should be performed[0m[0m
[34m[34m     [0m[0m [94m[94mfunction __fleek_prepare_completions[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug ""[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "========= starting completion logic =========="[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    # Start fresh[0m[0m
[34m[34m     [0m[0m [94m[94m    set --erase __fleek_comp_results[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_perform_completion_once[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "Completion results: $__fleek_perform_completion_once_result"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    if test -z "$__fleek_perform_completion_once_result"[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_debug "No completion, probably due to a failure"[0m[0m
[34m[34m     [0m[0m [94m[94m        # Might as well do file completion, in case it helps[0m[0m
[34m[34m     [0m[0m [94m[94m        return 1[0m[0m
[34m[34m     [0m[0m [94m[94m    end[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    set -l directive (string sub --start 2 $__fleek_perform_completion_once_result[-1])[0m[0m
[34m[34m     [0m[0m [94m[94m    set --global __fleek_comp_results $__fleek_perform_completion_once_result[1..-2][0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "Completions are: $__fleek_comp_results"[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "Directive is: $directive"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    set -l shellCompDirectiveError 1[0m[0m
[34m[34m     [0m[0m [94m[94m    set -l shellCompDirectiveNoSpace 2[0m[0m
[34m[34m     [0m[0m [94m[94m    set -l shellCompDirectiveNoFileComp 4[0m[0m
[34m[34m     [0m[0m [94m[94m    set -l shellCompDirectiveFilterFileExt 8[0m[0m
[34m[34m     [0m[0m [94m[94m    set -l shellCompDirectiveFilterDirs 16[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    if test -z "$directive"[0m[0m
[34m[34m     [0m[0m [94m[94m        set directive 0[0m[0m
[34m[34m     [0m[0m [94m[94m    end[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    set -l compErr (math (math --scale 0 $directive / $shellCompDirectiveError) % 2)[0m[0m
[34m[34m     [0m[0m [94m[94m    if test $compErr -eq 1[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_debug "Received error directive: aborting."[0m[0m
[34m[34m     [0m[0m [94m[94m        # Might as well do file completion, in case it helps[0m[0m
[34m[34m     [0m[0m [94m[94m        return 1[0m[0m
[34m[34m     [0m[0m [94m[94m    end[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    set -l filefilter (math (math --scale 0 $directive / $shellCompDirectiveFilterFileExt) % 2)[0m[0m
[34m[34m     [0m[0m [94m[94m    set -l dirfilter (math (math --scale 0 $directive / $shellCompDirectiveFilterDirs) % 2)[0m[0m
[34m[34m     [0m[0m [94m[94m    if test $filefilter -eq 1; or test $dirfilter -eq 1[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_debug "File extension filtering or directory filtering not supported"[0m[0m
[34m[34m     [0m[0m [94m[94m        # Do full file completion instead[0m[0m
[34m[34m     [0m[0m [94m[94m        return 1[0m[0m
[34m[34m     [0m[0m [94m[94m    end[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    set -l nospace (math (math --scale 0 $directive / $shellCompDirectiveNoSpace) % 2)[0m[0m
[34m[34m     [0m[0m [94m[94m    set -l nofiles (math (math --scale 0 $directive / $shellCompDirectiveNoFileComp) % 2)[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    __fleek_debug "nospace: $nospace, nofiles: $nofiles"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    # If we want to prevent a space, or if file completion is NOT disabled,[0m[0m
[34m[34m     [0m[0m [94m[94m    # we need to count the number of valid completions.[0m[0m
[34m[34m     [0m[0m [94m[94m    # To do so, we will filter on prefix as the completions we have received[0m[0m
[34m[34m     [0m[0m [94m[94m    # may not already be filtered so as to allow fish to match on different[0m[0m
[34m[34m     [0m[0m [94m[94m    # criteria than the prefix.[0m[0m
[34m[34m     [0m[0m [94m[94m    if test $nospace -ne 0; or test $nofiles -eq 0[0m[0m
[34m[34m     [0m[0m [94m[94m        set -l prefix (commandline -t | string escape --style=regex)[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_debug "prefix: $prefix"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m        set -l completions (string match -r -- "^$prefix.*" $__fleek_comp_results)[0m[0m
[34m[34m     [0m[0m [94m[94m        set --global __fleek_comp_results $completions[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_debug "Filtered completions are: $__fleek_comp_results"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m        # Important not to quote the variable for count to work[0m[0m
[34m[34m     [0m[0m [94m[94m        set -l numComps (count $__fleek_comp_results)[0m[0m
[34m[34m     [0m[0m [94m[94m        __fleek_debug "numComps: $numComps"[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m        if test $numComps -eq 1; and test $nospace -ne 0[0m[0m
[34m[34m     [0m[0m [94m[94m            # We must first split on \t to get rid of the descriptions to be[0m[0m
[34m[34m     [0m[0m [94m[94m            # able to check what the actual completion will be.[0m[0m
[34m[34m     [0m[0m [94m[94m            # We don't need descriptions anyway since there is only a single[0m[0m
[34m[34m     [0m[0m [94m[94m            # real completion which the shell will expand immediately.[0m[0m
[34m[34m     [0m[0m [94m[94m            set -l split (string split --max 1 \t $__fleek_comp_results[1])[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m            # Fish won't add a space if the completion ends with any[0m[0m
[34m[34m     [0m[0m [94m[94m            # of the following characters: @=/:.,[0m[0m
[34m[34m     [0m[0m [94m[94m            set -l lastChar (string sub -s -1 -- $split)[0m[0m
[34m[34m     [0m[0m [94m[94m            if not string match -r -q "[@=/:.,]" -- "$lastChar"[0m[0m
[34m[34m     [0m[0m [94m[94m                # In other cases, to support the "nospace" directive we trick the shell[0m[0m
[34m[34m     [0m[0m [94m[94m                # by outputting an extra, longer completion.[0m[0m
[34m[34m     [0m[0m [94m[94m                __fleek_debug "Adding second completion to perform nospace directive"[0m[0m
[34m[34m     [0m[0m [94m[94m                set --global __fleek_comp_results $split[1] $split[1].[0m[0m
[34m[34m     [0m[0m [94m[94m                __fleek_debug "Completions are now: $__fleek_comp_results"[0m[0m
[34m[34m     [0m[0m [94m[94m            end[0m[0m
[34m[34m     [0m[0m [94m[94m        end[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m        if test $numComps -eq 0; and test $nofiles -eq 0[0m[0m
[34m[34m     [0m[0m [94m[94m            # To be consistent with bash and zsh, we only trigger file[0m[0m
[34m[34m     [0m[0m [94m[94m            # completion when there are no other completions[0m[0m
[34m[34m     [0m[0m [94m[94m            __fleek_debug "Requesting file completion"[0m[0m
[34m[34m     [0m[0m [94m[94m            return 1[0m[0m
[34m[34m     [0m[0m [94m[94m        end[0m[0m
[34m[34m     [0m[0m [94m[94m    end[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m    return 0[0m[0m
[34m[34m     [0m[0m [94m[94mend[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m# Since Fish completions are only loaded once the user triggers them, we trigger them ourselves[0m[0m
[34m[34m     [0m[0m [94m[94m# so we can properly delete any completions provided by another script.[0m[0m
[34m[34m     [0m[0m [94m[94m# Only do this if the program can be found, or else fish may print some errors; besides,[0m[0m
[34m[34m     [0m[0m [94m[94m# the existing completions will only be loaded if the program can be found.[0m[0m
[34m[34m     [0m[0m [94m[94mif type -q "fleek"[0m[0m
[34m[34m     [0m[0m [94m[94m    # The space after the program name is essential to trigger completion for the program[0m[0m
[34m[34m     [0m[0m [94m[94m    # and not completion of the program name itself.[0m[0m
[34m[34m     [0m[0m [94m[94m    # Also, we use '> /dev/null 2>&1' since '&>' is not supported in older versions of fish.[0m[0m
[34m[34m     [0m[0m [94m[94m    complete --do-complete "fleek " > /dev/null 2>&1[0m[0m
[34m[34m     [0m[0m [94m[94mend[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m# Remove any pre-existing completions for the program since we will be handling all of them.[0m[0m
[34m[34m     [0m[0m [94m[94mcomplete -c fleek -e[0m[0m
[34m[34m     [0m[0m [94m[94m[0m[0m
[34m[34m     [0m[0m [94m[94m# this will get called after the two calls below and clear the $__fleek_perform_completion_once_result global[0m[0m
[34m[34m     [0m[0m [94m[94mcomplete -c fleek -n '__fleek_clear_perform_completion_once_result'[0m[0m
[34m[34m     [0m[0m [94m[94m# The call to __fleek_prepare_completions will setup __fleek_comp_results[0m[0m
[34m[34m     [0m[0m [94m[94m# which provides the program's completion choices.[0m[0m
[34m[34m     [0m[0m [94m[94m# If this doesn't require order preservation, we don't use the -k flag[0m[0m
[34m[34m     [0m[0m [94m[94mcomplete -c fleek -n 'not __fleek_requires_order_preservation && __fleek_prepare_completions' -f -a '$__fleek_comp_results'[0m[0m
[34m[34m     [0m[0m [94m[94m# otherwise we use the -k flag[0m[0m
[34m[34m     [0m[0m [94m[94mcomplete -k -c fleek -n '__fleek_requires_order_preservation && __fleek_prepare_completions' -f -a '$__fleek_comp_results'[0m[0m
