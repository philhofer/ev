#include "textflag.h"

TEXT ·__pipe2(SB),NOSPLIT,$0-24
	MOVQ $293, AX
	MOVQ dst+0(FP), DI
	MOVQ flags+8(FP), SI
	SYSCALL
	MOVQ AX, ret+16(FP)
	RET
	
TEXT ·epollcreate1(SB),NOSPLIT,$0-16
	MOVQ $291, AX
	MOVQ flags+0(FP), DI
	SYSCALL
	MOVQ AX, ret+8(FP)
	RET


TEXT ·epollwait(SB),NOSPLIT,$0-40
	CALL runtime·entersyscall(SB)
	MOVQ $232, AX
	MOVQ fd+0(FP), DI
	MOVQ events+8(FP), SI
	MOVQ nev+16(FP), DX
	MOVQ ms+24(FP), R10
	SYSCALL
	MOVQ AX, ret+32(FP)
	CALL runtime·exitsyscall(SB)
	RET


TEXT ·epollctl(SB),NOSPLIT,$0-40
	MOVQ $233, AX
	MOVQ fd+0(FP), DI
	MOVQ op+8(FP), SI
	MOVQ nfd+16(FP), DX
	MOVQ ev+24(FP), R10
	SYSCALL
	MOVQ AX, ret+32(FP)
	RET


TEXT ·writev(SB),NOSPLIT,$0-32
	CALL runtime·entersyscall(SB)
	MOVQ $20, AX
	MOVQ fd+0(FP), DI
	MOVQ iov+8(FP), SI
	MOVQ nev+16(FP), DX
	SYSCALL
	MOVQ AX, ret+24(FP)
	CALL runtime·exitsyscall(SB)
	RET


TEXT ·readv(SB),NOSPLIT,$0-32
	CALL runtime·entersyscall(SB)
	MOVQ $19, AX
	MOVQ fd+0(FP), DI
	MOVQ iov+8(FP), SI
	MOVQ nev+16(FP), DX
	SYSCALL
	MOVQ AX, ret+24(FP)
	CALL runtime·exitsyscall(SB)
	RET

