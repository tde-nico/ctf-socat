
NAME = socats

all: $(NAME)

$(NAME): main.go
	@ go build

run: all
	@ ./$(NAME) -c ./getflag -d

.PHONY: all, run
