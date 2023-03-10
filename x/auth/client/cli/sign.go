package cli

import (
	"fmt"
	"github.com/spf13/viper"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
)

const (
	flagAppend    = "append"
	flagPrintSigs = "print-sigs"
	flagOffline   = "offline"
)

// GetSignCommand returns the sign command
func GetSignCommand(codec *amino.Codec, decoder auth.AccountDecoder) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign <file>",
		Short: "Sign transactions generated offline",
		Long: `Sign transactions created with the --generate-only flag.
Read a transaction from <file>, sign it, and print its JSON encoding.

The --offline flag makes sure that the client will not reach out to the local cache.
Thus account number or sequence number lookups will not be performed and it is
recommended to set such parameters manually.`,
		RunE: makeSignCmd(codec, decoder),
		Args: cobra.ExactArgs(1),
	}
	cmd.Flags().String(client.FlagName, "", "Name of private key with which to sign")
	cmd.Flags().Bool(flagAppend, true, "Append the signature to the existing ones. If disabled, old signatures would be overwritten")
	cmd.Flags().Bool(flagPrintSigs, false, "Print the addresses that must sign the transaction and those who have already signed it, then exit")
	return cmd
}

func makeSignCmd(cdc *amino.Codec, decoder auth.AccountDecoder) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) (err error) {
		stdTx, err := readAndUnmarshalStdTx(cdc, args[0])
		if err != nil {
			return
		}

		if viper.GetBool(flagPrintSigs) {
			printSignatures(stdTx)
			return nil
		}

		name := viper.GetString(client.FlagName)
		cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(decoder)
		txBldr := authtxb.NewTxBuilderFromCLI()
		if len(txBldr.ChainID) == 0 {
			return fmt.Errorf("chain-id is missing")
		}

		newTx, err := utils.SignStdTx(txBldr, cliCtx, name, stdTx, viper.GetBool(flagAppend), viper.GetBool(flagOffline))
		if err != nil {
			return err
		}
		var json []byte
		if cliCtx.Indent {
			json, err = cdc.MarshalJSONIndent(newTx, "", "  ")
		} else {
			json, err = cdc.MarshalJSON(newTx)
		}
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", json)
		return
	}
}

func printSignatures(stdTx auth.StdTx) {
	fmt.Println("Signers:")
	for i, signer := range stdTx.GetSigners() {
		fmt.Printf(" %v: %v\n", i, signer.String())
	}
	fmt.Println("")
	fmt.Println("Signatures:")
	for i, sig := range stdTx.GetSignatures() {
		fmt.Printf(" %v: %v\n", i, sdk.AccAddress(sig.Address()).String())
	}
	return
}

func readAndUnmarshalStdTx(cdc *amino.Codec, filename string) (stdTx auth.StdTx, err error) {
	var bytes []byte
	if bytes, err = os.ReadFile(filename); err != nil {
		return
	}
	if err = cdc.UnmarshalJSON(bytes, &stdTx); err != nil {
		return
	}
	return
}
