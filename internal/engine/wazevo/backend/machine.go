package backend

import (
	"context"

	"github.com/tetratelabs/wazero/internal/engine/wazevo/backend/regalloc"
	"github.com/tetratelabs/wazero/internal/engine/wazevo/ssa"
	"github.com/tetratelabs/wazero/internal/engine/wazevo/wazevoapi"
)

type (
	// Machine is a backend for a specific ISA machine.
	Machine interface {
		ExecutableContext() ExecutableContext

		// DisableStackCheck disables the stack check for the current compilation for debugging/testing.
		DisableStackCheck()

		// InitializeABI initializes the FunctionABI for the given signature.
		InitializeABI(sig *ssa.Signature)

		// SetCompiler sets the compilation context used for the lifetime of Machine.
		// This is only called once per Machine, i.e. before the first compilation.
		SetCompiler(Compiler)

		// LowerSingleBranch is called when the compilation of the given single branch is started.
		LowerSingleBranch(b *ssa.Instruction)

		// LowerConditionalBranch is called when the compilation of the given conditional branch is started.
		LowerConditionalBranch(b *ssa.Instruction)

		// LowerInstr is called for each instruction in the given block except for the ones marked as already lowered
		// via Compiler.MarkLowered. The order is reverse, i.e. from the last instruction to the first one.
		//
		// Note: this can lower multiple instructions (which produce the inputs) at once whenever it's possible
		// for optimization.
		LowerInstr(*ssa.Instruction)

		// Reset resets the machine state for the next compilation.
		Reset()

		// InsertMove inserts a move instruction from src to dst whose type is typ.
		InsertMove(dst, src regalloc.VReg, typ ssa.Type)

		// InsertReturn inserts the return instruction to return from the current function.
		InsertReturn()

		// InsertLoadConstant inserts the instruction(s) to load the constant value into the given regalloc.VReg.
		InsertLoadConstant(instr *ssa.Instruction, vr regalloc.VReg)

		// Format returns the string representation of the currently compiled machine code.
		// This is only for testing purpose.
		Format() string

		// RegAlloc does the register allocation after lowering.
		RegAlloc()

		// SetupPrologue inserts the prologue after register allocations.
		SetupPrologue()

		// SetupEpilogue inserts the epilogue after register allocations.
		// This sets up the instructions for the inverse of SetupPrologue right before
		SetupEpilogue()

		// ResolveRelativeAddresses resolves the relative addresses after register allocations and prologue/epilogue setup.
		// After this, the compiler is finally ready to emit machine code.
		ResolveRelativeAddresses(ctx context.Context)

		// ResolveRelocations resolves the relocations after emitting machine code.
		ResolveRelocations(refToBinaryOffset map[ssa.FuncRef]int, binary []byte, relocations []RelocationInfo)

		// Encode encodes the machine instructions to the Compiler.
		Encode()

		// CompileGoFunctionTrampoline compiles the trampoline function  to call a Go function of the given exit code and signature.
		CompileGoFunctionTrampoline(exitCode wazevoapi.ExitCode, sig *ssa.Signature, needModuleContextPtr bool) []byte

		// CompileStackGrowCallSequence returns the sequence of instructions shared by all functions to
		// call the stack grow builtin function.
		CompileStackGrowCallSequence() []byte

		// CompileEntryPreamble returns the sequence of instructions shared by multiple functions to
		// enter the function from Go.
		CompileEntryPreamble(signature *ssa.Signature) []byte

		// LowerParams lowers the given parameters.
		LowerParams(params []ssa.Value)

		// LowerReturns lowers the given returns.
		LowerReturns(returns []ssa.Value)
	}
)
