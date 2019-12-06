#!/bin/bash

bail() {
    if [ -t 0 ]; then
        echo "$@" 1>&2
    else
        exec swaynag --edge bottom --message "$@"
    fi
    exit 1
}
dmenu() {
    dmenu_prompt=" toggle${DMENU_PROMPT:-}"
    if [ -x "$(command -v wofi)" ]; then
        exec wofi --cache-file /dev/null --allow-markup --dmenu --prompt "$dmenu_prompt"
    fi
    if [ -x "$(command -v rofi)" ]; then
        exec rofi -markup -dmenu -mesg "$dmenu_prompt"
    fi
    if [ -x "$(command -v dmenu)" ]; then
        exec dmenu -mesg "$dmenu_prompt"
    fi
    bail "no dmenu command found, tried wofi, rofi and dmenu."
}

ACTIONS=(
    "Start timer"
    "Stop timer"
    "List timers"
)
select_action() {
    for a in "${ACTIONS[@]}"; do echo "$a"; done | dmenu
}

select_project() {
    ret=$(toggl -t "$TOGGL_TOKEN" projects ls | DMENU_PROMPT=": <i>select a project</i>" dmenu | cut -d' ' -f 1)
    if [[ -z $ret ]]; then
        return
    fi
    toggl -t "$TOGGL_TOKEN" -j projects ls | jq -r ".[] | select(.id==$ret) | .name"
}

if [ -z "$TOGGL_TOKEN" ]; then
    bail "no toggl api token"
fi

case $(select_action) in
    "${ACTIONS[0]}") # start timer
        p=$(select_project)
        if [[ -z $p ]]; then
            exit 1
        fi
        d=$(DMENU_PROMPT=": <i>enter a description</i>" dmenu)
        ret=$(toggl -t "$TOGGL_TOKEN" timer start --description "$d" --project "$p")
        if [[ $? -ne 0 ]]; then
            bail "$ret"
        fi
        ;;
    "${ACTIONS[1]}") # stop timer
        ret=`toggl -t $TOGGL_TOKEN timers ls | tail -n +2 | cut -d' ' -f1`
        if [[ $? -ne 0 ]]; then
            bail "$ret"
        fi
        if [[ -z $ret ]]; then
            exit 1
        fi
        ret=$(toggl -t "$TOGGL_TOKEN" timer stop "$ret")
        if [[ $? -ne 0 ]]; then
            bail "$ret"
        fi
        ;;
    "${ACTIONS[2]}") # list timers
        toggl -t "$TOGGL_TOKEN" timers ls -a | dmenu
        ;;
esac
