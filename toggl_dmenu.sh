#!/bin/bash

bail() {
    exec rofi -e "$@"
    exit 1
}
rofi() {
    exec rofi -markup -markup-rows -dmenu -p toggl -i "$@"
}

ACTIONS=(
    "Start timer"
    "Stop timer"
    "List timers"
)
select_action() {
    for a in "${ACTIONS[@]}"; do echo "$a"; done | rofi
}

select_project() {
    ret=`toggl -t $TOKEN projects ls | rofi -mesg "<i>select a project</i>" -selected-row 1 | cut -d' ' -f 1`
    if [[ -z $ret ]]; then
        return
    fi
    toggl -t $TOKEN -j projects ls | jq -r ".[] | select(.id==$ret) | .name"
}

if [[ -z $TOKEN ]]; then
    bail "no toggl api token"
fi

case `select_action` in
    "${ACTIONS[0]}") # start timer
        p=`select_project`
        if [[ -z $p ]]; then
            exit 1
        fi
        d=`rofi -mesg "<i>enter a description</i>"`
        ret=`toggl -t $TOKEN timer start --description "$d" --project "$p"`
        if [[ $? -ne 0 ]]; then
            bail "$ret"
        fi
        ;;
    "${ACTIONS[1]}") # stop timer
        ret=`toggl -t $TOKEN timers ls | tail -n +2 | cut -d' ' -f1`
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
    "${ACTIONS[2]}") # running timer
        toggl -t $TOKEN timers ls -a | rofi -selected-row 1
        ;;
esac
