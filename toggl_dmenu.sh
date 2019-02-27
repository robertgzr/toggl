#!/bin/bash

bail() {
    exec rofi -e "$@"
    exit 1
}
rofi() {
    exec rofi -dmenu -p toggl -i "$@"
}

ACTIONS=(
    "Start timer"
    "Stop timer"
    "Running timers"
)
select_action() {
    for a in "${ACTIONS[@]}"; do echo "$a"; done | rofi
}

select_project() {
    ret=`toggl -t $TOKEN projects ls | rofi -mesg "select a project" -selected-row 1 | cut -d' ' -f 1`
    if [[ -z $ret ]]; then
        return
    fi
    toggl -t $TOKEN -j projects ls | jq -r ".[] | select(.id==$ret) | .name"
}

if [[ -z $TOKEN ]]; then
    bail "no toggl api token"
fi

case `select_action` in
    "${ACTIONS[0]}")
        p=`select_project`
        if [[ -z $p ]]; then
            exit 1
        fi
        d=`rofi -mesg "enter a description"`
        if [[ -z $d ]]; then
            exit 1
        fi
        ret=`toggl -t $TOKEN timer start --description "$d" --project "$p"`
        if [[ $? -ne 0 ]]; then
            bail "$ret"
        fi
        ;;
    "${ACTIONS[1]}")
        ret=`toggl -t $TOKEN timers ls | rofi -selected-row 1 | cut -d' ' -f1`
        if [[ $? -ne 0 ]]; then
            bail "$ret"
        fi
        if [[ -z $ret ]]; then
            exit 1
        fi
        ret=`toggl -t $TOKEN timer stop "$ret"`
        if [[ $? -ne 0 ]]; then
            bail "$ret"
        fi
        ;;
    "${ACTIONS[2]}")
        toggl -t $TOKEN timers ls | rofi -selected-row 1
        ;;
esac
