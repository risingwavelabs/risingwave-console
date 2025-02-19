#!/usr/bin/awk -f

# This script performs template substitution by reading commands and their corresponding template variables
# Usage: awk -f template-subst.awk -v cmds='cmd1|var1;cmd2|var2' template_file > output_file

BEGIN {
    split(cmds, cmd_pairs, ";")
    for (i in cmd_pairs) {
        split(cmd_pairs[i], pair, "|")
        cmd = pair[1]
        var = pair[2]
        content = ""
        while ((cmd | getline line) > 0) {
            content = content line "\n"
        }
        close(cmd)
        templates[var] = content
    }
}

{
    line = $0
    for (var in templates) {
        gsub("{{" var "}}", templates[var], line)
    }
    print line
} 