#!/bin/bash

rsync -av --ignore-existing annotations/ deploy@citybits.co.uk:/srv/apps/citybits/notes/sb/annotations/