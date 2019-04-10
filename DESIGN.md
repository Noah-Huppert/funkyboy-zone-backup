# Design
Mountain backup design.

# Table Of Contents
- [Overview](#overview)

# Overview
The tool provides different modules which backup files in different ways.  

Each module implements the `Backuper` interface. This interface writes data to 
backup into a tar.gz file. The program then uploads this file after all backup 
modules have run.
