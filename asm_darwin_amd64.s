
#include "textflag.h"

TEXT ·kqueue(SB),NOSPLIT,$0-4
	MOVL $(0x2000000+362), AX
	SYSCALL
	JCC  2(PC)
	NEGQ AX
	MOVL AX, ret+0(FP)
	RET

TEXT ·kevent(SB),NOSPLIT,$0-52
	CALL runtime·entersyscall(SB)
	MOVL $(0x2000000+363), AX
	MOVQ fd+0(FP), DI
	MOVQ ch+8(FP), SI
	MOVQ nch+16(FP), DX
	MOVQ ev+24(FP), R10
	MOVQ nev+32(FP), R8
	MOVQ ts+40(FP), R9
	SYSCALL
	JCC  2(PC)
	NEGQ AX
	MOVL AX, ret+48(FP)
	CALL runtime·exitsyscall(SB)
	RET

// fcntl(fd, F_SETFD, FD_CLOEXEC)
TEXT ·cloexec(SB),NOSPLIT,$0-8
	MOVL $(0x2000000+92), AX
	MOVQ fd+0(FP), DI
	MOVQ $2, SI
	MOVQ $1, DX
	SYSCALL
	RET

TEXT ·__pipe(SB),NOSPLIT,$0-16
	MOVL $(0x2000000+42), AX
	SYSCALL
	JCC  2(PC)
	NEGQ AX
	MOVQ AX, ret+0(FP)
	MOVQ DX, ret+8(FP)
	RET 
