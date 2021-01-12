# Google Hash Code 2019

<a title="AmrAtrash / CC BY-SA (https://creativecommons.org/licenses/by-sa/4.0)"><img width="512" src="https://upload.wikimedia.org/wikipedia/commons/7/7e/Event_google-hash-code_491696.jpg"></a>

Hash Code is a coding competition organized by Google where teams from all over the world compete in solving exciting engineering problems.

This is a [Google Hash Code 2019](https://codingcompetitions.withgoogle.com/hashcode) solution: [Problem statement](https://storage.googleapis.com/coding-competitions.appspot.com/HC/2019/hashcode2019_qualification_task.pdf).

Disclaimer: The goal of this project was not necessarily to achieve the highest score possible, but to get to know golang and play around with its features.
Our solution achieves a score of around 509400, which would be equivalent to about the 1470th place out of 6640 participants.

## Build & Run

```bash
git clone https://github.com/baurst/go_hc19
cd go_hc19
go run .
```

By default, the solver will run on all of the following datasets in parallel:

* a_example.txt
* b_lovely_landscapes.txt
* c_memorable_moments.txt
* d_pet_pictures.txt
* e_shiny_selfies.txt

If you would like to solve a subset of the problems, use the flag --datasets for example like so: ```--datasets b_lovely_landscapes.txt```.
The tool will create a directory called out and save each solution in a seperate file, reporting the respective score.

## Highly subjective takeaways about doing coding competitions with Golang

* due to Golangs small feature set and explicit style, you end up writing slightly more code than you would in Python, C++ or Rust
* this makes good preparation essential: Have your boilerplate ready for data loading, writing solutions and asynchronous problem solving etc.
* the lack of functional programming (map, filter and reduce) negatively impacts coding speed
* the ease with which you can run goroutines asynchronously is just great

## Team Members

* [Jasmin Kling](https://github.com/jkling2)
* [Stefan Baur](https://github.com/baurst)
