#include "textflag.h"

TEXT ·semacquire(SB),NOSPLIT,$0-8
	MOVB $1, true+9(FP)
	JMP runtime·semacquire(SB)

TEXT ·semrelease(SB),NOSPLIT,$0-8
	JMP runtime·semrelease(SB)
