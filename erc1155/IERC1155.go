package erc1155

import "github.com/studyzy/openzeppelin-go/common"

/**
 * @dev Required interface of an ERC1155 compliant contract, as defined in the
 * https://eips.ethereum.org/EIPS/eip-1155[EIP].
 *
 * _Available since v3.1._
 */
type IERC1155 interface {
	SupportsInterface(interfaceId string) bool

	/**
	 * @dev Returns the amount of tokens of token type `id` owned by `account`.
	 *
	 * Requirements:
	 *
	 * - `account` cannot be the zero address.
	 */
	BalanceOf(account common.Account, id *common.SafeUint256) (*common.SafeUint256, error)

	/**
	 * @dev xref:ROOT:erc1155.adoc#batch-operations[Batched] version of {balanceOf}.
	 *
	 * Requirements:
	 *
	 * - `accounts` and `ids` must have the same length.
	 */
	BalanceOfBatch(accounts []common.Account, ids []*common.SafeUint256) ([]*common.SafeUint256, error)

	/**
	 * @dev Grants or revokes permission to `operator` to transfer the caller's tokens, according to `approved`,
	 *
	 * Emits an {ApprovalForAll} event.
	 *
	 * Requirements:
	 *
	 * - `operator` cannot be the caller.
	 */
	SetApprovalForAll(operator common.Account, approved bool) error

	/**
	 * @dev Returns true if `operator` is approved to transfer ``account``'s tokens.
	 *
	 * See {setApprovalForAll}.
	 */
	IsApprovedForAll(account common.Account, operator common.Account) (bool, error)

	/**
	 * @dev Transfers `amount` tokens of token type `id` from `from` to `to`.
	 *
	 * Emits a {TransferSingle} event.
	 *
	 * Requirements:
	 *
	 * - `to` cannot be the zero address.
	 * - If the caller is not `from`, it must have been approved to spend ``from``'s tokens via {setApprovalForAll}.
	 * - `from` must have a balance of tokens of type `id` of at least `amount`.
	 * - If `to` refers to a smart contract, it must implement {IERC1155Receiver-onERC1155Received} and return the
	 * acceptance magic value.
	 */
	SafeTransferFrom(from, to common.Account, id, amount *common.SafeUint256, data []byte) error

	/**
	 * @dev xref:ROOT:erc1155.adoc#batch-operations[Batched] version of {safeTransferFrom}.
	 *
	 * Emits a {TransferBatch} event.
	 *
	 * Requirements:
	 *
	 * - `ids` and `amounts` must have the same length.
	 * - If `to` refers to a smart contract, it must implement {IERC1155Receiver-onERC1155BatchReceived} and return the
	 * acceptance magic value.
	 */
	SafeBatchTransferFrom(from, to common.Account, ids, amounts []*common.SafeUint256, data []byte) error
	/**
	 * @dev Returns the URI for token type `id`.
	 *
	 * If the `\{id\}` substring is present in the URI, it must be replaced by
	 * clients with the actual token type ID.
	 */
	Uri(id *common.SafeUint256) (string, error)
}
