#!/usr/bin/python3
# -*- coding: utf-8 -*-

# Usage: go run . samples/example02 | python3 scripts/moves.py

from sys import stdin
from re import search, MULTILINE, IGNORECASE
from collections import defaultdict


def extract_ant_move(move: str):
    match = search(r"L(?P<ant>\d+)-(?P<to>\w+)", move)
    return int(match.group('ant')), match.group('to')


def print_moves(ants):
    for ant_id, moves in ants.items():
        moves = " > ".join([f"T{t + 1}#{m}" for (t, m) in moves])
        print(f"{ant_id:<4} :: {moves}")


text = stdin.read()
match = search(r'Parsed.+$([\w\s.µ-]*)Finished', text, MULTILINE | IGNORECASE)
if match:
    ants = defaultdict(list)
    text = match.group(1).strip()
    for turn, turn_moves in enumerate(text.splitlines()):
        turn_moves = turn_moves.split()
        for move in turn_moves:
            ant_id, ant_moved_to = extract_ant_move(move)
            ants[ant_id].append((turn, ant_moved_to))
    print_moves(ants)
