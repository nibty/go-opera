// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package evmcore

import (
	"crypto/ecdsa"
	"errors"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"

	"github.com/Fantom-foundation/go-opera/inter"
)

var FakeGenesisTime = inter.Timestamp(1608600000 * time.Second)

// ApplyFakeGenesis writes or updates the genesis block in db.
func ApplyFakeGenesis(statedb *state.StateDB, time inter.Timestamp, balances map[common.Address]*big.Int) (*EvmBlock, error) {
	for acc, balance := range balances {
		statedb.SetBalance(acc, balance)
	}

	// initial block
	root, err := flush(statedb, true)
	if err != nil {
		return nil, err
	}
	block := genesisBlock(time, root)

	return block, nil
}

func flush(statedb *state.StateDB, clean bool) (root common.Hash, err error) {
	root, err = statedb.Commit(clean)
	if err != nil {
		return
	}
	err = statedb.Database().TrieDB().Commit(root, false, nil)
	if err != nil {
		return
	}

	if !clean {
		err = statedb.Database().TrieDB().Cap(0)
	}

	return
}

// genesisBlock makes genesis block with pretty hash.
func genesisBlock(time inter.Timestamp, root common.Hash) *EvmBlock {
	block := &EvmBlock{
		EvmHeader: EvmHeader{
			Number:   big.NewInt(0),
			Time:     time,
			GasLimit: math.MaxUint64,
			Root:     root,
			TxHash:   types.EmptyRootHash,
		},
	}

	return block
}

// MustApplyFakeGenesis writes the genesis block and state to db, panicking on error.
func MustApplyFakeGenesis(statedb *state.StateDB, time inter.Timestamp, balances map[common.Address]*big.Int) *EvmBlock {
	block, err := ApplyFakeGenesis(statedb, time, balances)
	if err != nil {
		log.Crit("ApplyFakeGenesis", "err", err)
	}
	return block
}

// FakeKey gets n-th fake private key.
func FakeKey(n uint32) *ecdsa.PrivateKey {
	var keys = [400]string{
		"0x163f5f0f9a621d72fedd85ffca3d08d131ab4e812181e0d30ffd1c885d20aac7",
		"0x3144c0aa4ced56dc15c79b045bc5559a5ac9363d98db6df321fe3847a103740f",
		"0x04a531f967898df5dbe223b67989b248e23c1c356a3f6717775cccb7fe53482c",
		"0x00ca81d4fe11c23fae8b5e5b06f9fe952c99ca46abaec8bda70a678cd0314dde",
		"0x532d9b2ce282fad94efefcf076fdfbe5befe558b145f4cc97f953bcabf087aeb",
		"0x6e50dbd3e81b22424cb230133b87bc9ef0f17c584a2a5dc4b212d2b83b5ee084",
		"0x2215aaee06a2d64ca32b201e1fb9d1e3c7a25d45a6d8b0de6300ba3a20e42ef5",
		"0x1cd6fdfc633c0fa73bd306c46eecd23096365b44ab75f0e6fa04dc2adbea9583",
		"0x2fc91d5829f44650c32ba92c8b29d511511446b91badf03b1fd0f808b91a4b5b",
		"0x6aeeb7f09e757baa9d3935a042c3d0d46a2eda19e9b676283dce4eaf32e29dc9",
		"0x7d51a817ee07c3f28581c47a5072142193337fdca4d7911e58c5af2d03895d1a",
		"0x59963733b8a6fb1c6eeb1ce51c7e6046e652a9bcacd4cbaa3f6f26dafe7f79f7",
		"0x4cf757812428b0764a871e94b02ba026a5d3738e69f7d1d4f9f93b43ed00e820",
		"0xa80a59dc6a9be8003a696ed08a4d37d5046f66201912b40c224d4fe96b515231",
		"0xa2ef6534312d205b045a94ec2e9d49191a6d17702671d51dd88a9e2837b612ce",
		"0x9512765baac04484c19491feb59fe8ef8ba557e29e00100f3159c8ee35c89038",
		"0x700717777c4b7ccdc8c79d6823cb3ea0356ac5e3822accdfa8539cf833caae15",
		"0x838ab204f288b4673bbc603ac52c167e8b1c1392cdd96bc02b8fcfadec98cc26",
		"0xbf6ba360590e69d1495ea8c0ab2f4a18ebbed7c4bbbe2d823a57719cd40df94f",
		"0xd2f091785e9ca0ea2388cf90a046e87552e5cbb4492a9702b83aa32dddf821ac",
		"0xad8b51bf6a35a934587413394ab453df559603f885ae3ff0e90c1d90c78153bd",
		"0xa1ae301b83bfdd9e4a6ffb277896e5b4438725844fd44b5f733219f9f0a1402b",
		"0x9bf39f28aa39153777677711b8ca8a733ffcdad9ec8831713a01d71fb3dbe184",
		"0xee948e4413ce4e82ecb51fb6669f82d5af9b0ca4c31924514e6e844e8da46051",
		"0xe9a94ddcec56059cffb6dd699011f2bb323293f90613385c8624839296b3d182",
		"0xbccc8d4364e82a04ea2dc840ad6eeec6a2c35a51fb01943d58728da7bd4364dc",
		"0xd8af1e1f98a3628e91e46888b02cb34b00fd72aee1946409a3435ea806f1ace8",
		"0x0ba54f6d7c269ae7d115a17446abe7ba52293997de821d262a3d113fd694d85a",
		"0x2666d00809c1ce11da2c7598d3ab54e1bb75263d9e25d8209568a1d5e7cf9cb7",
		"0x204b21603d4a076bcdd34db298229f935198ea695964a3e156289728290e6240",
		"0xd30f524750dbbef5833dddbcdcdaf6c7f4e43c777d5c468d124170838e83c59d",
		"0xee3002a37a510360b0de793b45dd56a4a7a1df843e04cd991441854978b5154e",
		"0x410605310cdc3bb8ffe3d472ebf183545e6a09f3b211616156d42d8ad2ee1218",
		"0x3d47844c536f73c3558bf2e2238b13b327be2890cad6de60a3940337b8afa774",
		"0x114a976408f9a71c581871a6b68f5006df44a178e86d0bc7659d591bb4e56da6",
		"0x2c1108cc823ae0c5f496ad61eaab90d0677875ab1b2e0a55a89ecd87388fa9b3",
		"0x5fb2462733c28810a8bc68712a08c201ed2b89e822bc7309834476cfa1857acc",
		"0x595726750f55bc28a9a2e50f92a6c5fab42e738409cab0008299039c9966e0fe",
		"0x6cb89990a3ecf4930470351f1d76a72525d2f702e47d490cc0cc087893d2664a",
		"0x275b9ee8df6f2c2d02cd1fb5c343f139867104d5da6f8d82afc771e2d11a28e4",
		"0x3a5abe2f6ee961774f0466fca8f986fb0a53c5560b0f239d2a7ce0c8cdb3e1d1",
		"0x97bb9f4bb472e89fc5405dd5786ea1de87c5d79758020abb0bfcbf4c48daf9a2",
		"0x8ae00c99180aa1f9a8fe2ee7e53aaaedc0e55ff71785f26afa725295d4de78ff",
		"0x65e35bf4a097d543df427ec67c192765f6edcbdda54e1b5c0bb5e99507f6a269",
		"0x5fa4c34c434e0ddf1f8b3917c3e9ecfcbc09df91237012ff8105bcba275e4b7a",
		"0x52ebc273f1da45483d5c6d884f0b65dda891ffee0ea6cdb0c6af90e852984e96",
		"0xadec518fdc716a50ffc89c1caf6d2861ffaf02f06656d9275bd5e5d25696c840",
		"0xc08f211d4803a2ab2c6c3c0745671decba5564dbebf9d339ae3344d055fd8e1d",
		"0x7cf47f78fc8a5a550eae7bc77fb2943dbf8b37dfc381b35ccc6970684ac6cbee",
		"0x90659790bafc92adea25095ebfabaffaa5c4bf1d1cc375ab3eac825832181398",
		"0xebcae9b7ee8dc6b813fd7aa567f94d9986a7d39a4997ebea3b09db85941cedb5",
		"0xde2da353e4200f22614ce98b03e7af8e3f79afa4dcd40666b679525103301606",
		"0xd850eca0a7ac46bc01962bcff3cd27fff5e32d401a4a4db3883a3f0e0bdf0933",
		"0xabd5c47c0be87b046603c7016258e65c4345aab4a65dde1daf60f2fb3f6c7b0c",
		"0xa6c6c5d4c176336017fe715861750fe674b6552755010bd1e8f24cbee19b9b59",
		"0xf90b1f7c5e046b894c4c70d93ed74367c4ec260a5ee0051876c929b0a7e98dcb",
		"0x15316697d6979fd22c5e3127e030484d94e4e7d78419200e669c3c2d0d9aa2e4",
		"0xe86120a57411c86be886b6df0e80ee2061ddf322043ef565b47b548c8072ae31",
		"0xe3e66700a59d00d5778c6b732d0e5f90b1516881a76ee4ad232aa7d06ba11e62",
		"0xdd876a98e12334ef52ca2d8e149a20b5e085e7e8c6281c2aa60736915073f53f",
		"0x30ab7160e3c2ec3884117f91e5189ca1c16af03af36a75cc0169b5f2e8163a88",
		"0x2b2a8abbbc4624e33f737bbcf8d864999619e7eb2e92630c2ce3a773c438fba5",
		"0xffcad1487800293d08dbe6355f60c170f41ae93906293f2a30c00568f6fb8717",
		"0x1b70a964ce916046ba1d3def8fc7d004f213028113e2751e3cf0a12307a21e9f",
		"0x4eb12e7c8a1d99a00dc99df7f8c162a929894ab2a638048627a08d9913c02efd",
		"0x6944e7e33cafcd099c1b2a88e87e8f57b3fc48a0002c4d168737f55bca9dce6e",
		"0x3c2bf03d642d85932ef2f6cc23259f8cc8782c60043c9d7ae58b096a02f9007b",
		"0x161c258dc7afadecfe8f8106ab619690ac01f52f669a3b1f453540bf82c78b14",
		"0x2941eea6ed3ee2166a0a8ce17f4a7e571cd8fd23ca270cc72839d7bafa955845",
		"0x85449300aed219707b7801669597c082dd9e4c74633472610c0009e79422da53",
		"0x79254268cf6352f9910405bb8c545aeeee8fbb61293e62663f81355b0ba3d86b",
		"0x548c31c260958764c20b417b416223ab8623e8364d84fc8806f665eeb084d6d9",
		"0x67edd0d0682a4e7e52575884db12d03b325bd6f8a5e18fb65b143f9f25df04aa",
		"0xa262f0ea06bc87c7829dc6fe38d83f99c326976a13c5a0e94da13adc5d136307",
		"0x9c655e84994bbe21ac8bade72dd6ebe37491c2986a8e4eb6ab2d007e3f130270",
		"0xaf06f1eb84c1f3c9e14e35e6aed34703803fea21dcf628f2bc178325b33148c1",
		"0x6aad2dc6d0442b14cfbdcc0ae207c3b88d31e1294606c287d1a6b0cc58670a5f",
		"0x7e4e0e156ae1532e98bfda00d850fb476f43b81f09af080446fece7b7d8ac388",
		"0xdaf1bfaeb6700798569f6d4815dff9fa6590856c27e6d9aa112cb06d5f21b525",
		"0xcd76b0089ef74032756bf06fdd5903c8787957e943e0622f9b35ba84185cb675",
		"0xa7b791cf6aae777a7954961b6b5c66b9326ee35ce379f17f554a55bd7cc91d5f",
		"0xbb76f0697e3baca59a9b7e5ad449c912ad568d371a1c579bc92b8043479606fc",
		"0x95f963b910d735ffd9dfd784ca683cc790aa40acdc5fa9d27fe8e84b6900528d",
		"0xe8fa32934f5e6c895f1586114cd3c84caf04f01efab47bcd8bc2efbd743e6dfb",
		"0x030076ec9dd721e30c743bea5f0019e80428a97fe900b7245b968c6a6f774313",
		"0xb7001594545d584baab7d6b056d3a39c325db2b450787329ed916dd26fa32260",
		"0xd2a748ae6cfa9156c94754320b22b35b295dfe779887098386b2cea72f4d0dd2",
		"0x258bb5bc5a87c5aebdadb7f83f39029242f0650bc52fc4cc482c89412532db2b",
		"0x20ecd5168010ed59a184fd5b5ac528b9fff70ed103ea58c6f35f0137854969b8",
		"0x1a517bb0a54425b29b032d4cc28423224963cda64b50ef9731b429660c8da129",
		"0xee61ccd710d6dabb81849a1ca2a9369c5484b65517ea80c7772135c55aa9f147",
		"0xe906f9f17d6511279053a95feca4a55ef59f97f2596790163ade48067d20238f",
		"0x3c8cd94010f44ac08efaae8a31e726bb4fc95564ee7c2e03e7a43ee43f31d6ed",
		"0x58ed1a9bde738f0a089ae0c7e18bf77d2489e75f282c1af6e16e0b86fc30c41e",
		"0x0be230f443fdc623254292438f0c0f86ba0b799c89d15c9383aa3a64d99628e7",
		"0x26a53f1c2b8fff8cfc6ca691777d1c84eb6bef60a3b3c26233c1b12faf653584",
		"0x79a6be35f01db3d6a2b659b49f04505a0ce1abf3770494a7a24a83642ae8a675",
		"0x54cd7ec556aeeb6f538fe1b7523a86425e520e288be516c95909582e79012bc3",
		"0x675e4d1f766213db0791a274e748e858d091460fb4a9222e4b75380f8ecabcdc",
		"0x4203c438d4e94bf4a595794b5f5c2882f959face730abb7a7b8acb462c8e138d",
		"0xf1d724ebfb5f7b17f7052cdf7a836d4c4c7ee4dd3e7ca4725575484f888433e9",
		"0x09a90166b0a2ed9424133272e91a07c3360726ac80d6e1573060351f39480587",
		"0x5a46eb6c75bffb43b5fa663489020f3b6d4536085594ec8f19f337a3736ba6cc",
		"0xb1e7cf7d8a29debe76ec7cb70b41f8ac6469f3e2c8b37a2f14790d163d15370a",
		"0x1a52df53475ced5eb4e64c813eea27e936b47ff0b432c1e3d644277fe863c3a9",
		"0x553c9973f1fbb945a736a6bdfd18ba9d91fc06565eac586fa6588e09853d1df0",
		"0xfb95936263538230ed917f11f33ac41968edb54245c945cd54b2e9256a08a844",
		"0x556c3cde12eb0892a8b7c6b342981e2c254958420a685579a2adfd1142576d46",
		"0x0090df7ebb848c5972f46dbf6242a17f83287da350d32a395c3ab83058ba691b",
		"0x25955bdbd28945f50d17724bbced5af21151626af51592ffa994f10aa6633b4e",
		"0x68f75a048cd1ce89b393c19632116934315108953eb847a1f037d78ce6960666",
		"0xca91899d0a625d34c196eb3fffe50e51ae3e87aacab82f0c53e37520c7e45af3",
		"0x3f5a5e5bc69137653eabb3f5c90ae596196ccc9c9925e53ecdc853b56eeb784c",
		"0xde93ecd614da95a47bb88b5f9110a9a5a9f016bc9e88134bcfe324378342e5e9",
		"0x1a80db1ae69eb985e62d611c4d1493ed3de68e58931db08486e6c5283765514f",
		"0x9cbdcb0f85e566da2d6c64b856c8aa20e78ee0465e17c29d319d4a53ea2541ee",
		"0x470e3980e35910994c2bf0a7e64b59a28ccbc1d995ba634d76cc33f6a73004b3",
		"0x071e22b5fe73f34462eca438e49c34d85f09e3eb13c9f66dd5fed6d9689bbaf8",
		"0x139d795138f400fca79d909bbad6b035d7b947df7cdd50aaf9a616a24de8c130",
		"0x6e30409a92cc0c8cf86f8ae6ade4a7ce53bf625169bdfa02417de7627446f326",
		"0x915889a946129f582197a1db0a09a6f12d3b6c9efa6923909879d180604677e0",
		"0x0c82284d2fa12ed5fb43d9dae8f3debdf5cf26b8741ace39b135d0040a151be8",
		"0x9f0b6818feedea2c15d62da3f6dd48092c7d0e00559643c867bf071e5a9ff682",
		"0xb34666769f0451291b2d60d0e8f1001318c9f43adfd603052b3f381a5678defe",
		"0x4e2c4db08ae9c20f6126a8a0edd30968c2c55ce333f84813fe672dc9dffe5a6a",
		"0x7bdaae56660f7930d78eaf717381a207e4e86985179ffc6914446823127a7c16",
		"0xff0da7ec6e3e56fadf382daa3d611fe1328c5219d4c2fb73c09a8fcb12bcf91b",
		"0xdb04d9f9d8c937b02c8d467ef377b47abe4516d1eb32b78712ec45e5516c1ba7",
		"0x2450c3cf1574a70e78a82af4a1abb6382bf0e7e1296997f2315139fcc595e01a",
		"0xb41c0ec277709073486a28f03a31709848c0d57037b8bd7a7debbc178f90115a",
		"0x7c43ac7fa076420af6e8645b3200d5d57350ff77cd911af371e5b653a6dd6f8e",
		"0xf8ed2daf53f7243a508ee65ccd195f7dde1999803d807ff15ec51fb04d02613d",
		"0xe606d3daecd9f06fd60131f9002d3491c9e86a08adf56323dce8e7f06e095319",
		"0x2974a674b028fc36315ceac006cc6fd858ea45e4b53176cf092f5c40560fd508",
		"0x8a35142859ad80c41626584538b24405f3ae2733cc41147662509ce923bb83b9",
		"0x2abd57cf597c98f22abca39756430da944782ddee87380cc56c11815e5eaf47c",
		"0x7450b15b4dca9639a538b6403fe8c99aa5de8244cf5064c89a85ccd4dcb573da",
		"0x486075ba7a2cdeb60b306d1c8a3680831b7f1634796acac1189a4b3995509c9a",
		"0xd84da5955eac7f6a0a6daf29e69b236837ff3de881d0626f4bb1e0f2ffbced6d",
		"0x5af34614f968625ecb3036367997258f13c3074c696fa5e74b9a20f6c9e6d6e8",
		"0x7f2718d8bdb921b2b8e3b8f580df3509727d0e6c416241080a8f07784344ecfd",
		"0x2fe9624cd4e14a5d090048348b606dd1f923f005fae6a3fe49ca97eba0361d53",
		"0x3f2fc9c73404966f095ef4a4b720a5e6bc66588d919a1d93f61ef202f212917f",
		"0x71f70b4dee2fa168c60d975d4c099e6afb4111d4a25a7e6060d6ecf91a762ff2",
		"0x1f536c448981894a41f77a854c7c87813460b95e04f1721e0c22c525058f618d",
		"0x7d2acd17fe73db065d71bdab05e31c3980174f1866021978ab6a88c312798b02",
		"0x6f579d82fe9bbf69dc0948a1e70fbb492853735d2c64ae96e02630fd5b362750",
		"0x174fdd2d55d3a000c3c2089f4a9f0fe44bc529147e5c0f365ead8cf6793cd875",
		"0x294a5ea034eda467a892984045ba74babca1fef909f963153e03a216d7125155",
		"0x07ed10820b5f5f16ff65f93f889ec65ce609761873f4dec55f8283ad7505bd17",
		"0x7d0bae5c434f420d55a635ae7909167bf326a6de119fe3bd52806142bb0a6c9a",
		"0xe5110b52da5c9bb29077b3871f5337a182a434834f08aa0de9ce2d3f84044441",
		"0x5efa44aa50b6263cb442ce41b0a99f80354bd1bfd5bde0448844be7ab37e7321",
		"0x662595fd57f899dbfa3b1dc69eb61f6a0be7df0c646d7fb2f1fb19dccf61ca17",
		"0x17a25d6fa0c24d573f0251c666ec730e5046f061bca56a69c350e710b9a6d0ed",
		"0x5d6fe06f8b7fdbfb1b2f01576ed4901767fadd740a5b4ccb72ec0ad9b35e29d1",
		"0x7da99f9a2621bb6f330c7ca4040e1fd009c71edaf7aa0937aba5a3d0974abe66",
		"0xd81aa645822d0cf9d4d05e00bf8f2d098e1f6f8d32b8cb687d64eff5498c11cf",
		"0x5d55d5313ab4768eb512510d02bd3a23e8155821f86edfaf0b74d8596554b9fb",
		"0x127688999f47d75aba54d8aaeb180fd436e0da0e43377a7b8d69f19eb4578a99",
		"0x49c37d4adca011452b1d57fb6909bbca8fea4c21024266a288169f0afc4b51c7",
		"0x0829b7097e64c7d36432585998e07e3ffa3efe7fab14bcd65d7ab57fc368cef4",
		"0xc66691853c8e8b39081307b5f3f2dc7e3093dfa16675d7918ee4375a6a438c7f",
		"0xa5d2f7d5653ca14a7a6aeafd72a8616c67c1dc06a8d1b542441e1a8c4c66b698",
		"0x0d29226859dfc67e625524722a710010078dcee30f419e9135cf37e5f62ad992",
		"0x093dee0ec967dc27be52ef2a5075de189fb86c1d9ac8a323a21ed28322c484cf",
		"0x6d251a4c5c4061c45bbd84df341fe13c38781a376477c3faf87ec01affbc17da",
		"0x2ce0089d7f7fa183dfdbddd0f9b955515afab0da8200b9b3f54f98c0e1c508e3",
		"0x12ccde5e31fb4a95e7c4f772b1327df6478c3eb418a92d71c83a022a7e5befdf",
		"0x2e95b6e06eb044930de208688405d124d6128441cce19f26a85bbecf3e535f59",
		"0x9f1e600123f94ace5ce28be6162b4aba7d553e38528f730467842bdf540d9ae0",
		"0x4e495029d169549625d582fd01e19c04cc09a73fda551815b667156c3498d952",
		"0xb49fc99bcc07067f537986c573c4d86d338461cb46180cedf3419bc64bd245b7",
		"0xbfe59e3c21596612388930a1ce40a41e6a4e9bd66cdf8835f3c6f2372b319900",
		"0x80596ce3db698182addf4955e353c04915f6bffacebd1165fac1daeaada1bed7",
		"0x56fb2cf7d4213cf2a52d1aa3cf7b8accdd7be4c4c12a27b52e1eac576b3d750d",
		"0xd73ff73d7f2ee5d1d918ebbbae3281b4af586b20d57de01fa452228d07759f89",
		"0xbe751235776e44b6b2ca4c5a9e0c1d4313e87fcf24c1dec784303dc674dac2d3",
		"0x3a8c195c8ffb62a6aabf57f1f37841548d1251368a7e891d04c758aae1bdedea",
		"0xc8a2e605347e0585ba927a265c3624e60687f0d2e275301ed21a0b0422160cb1",
		"0xf055778b054a2c38010fbc2848d6c62e9dee1966251bbfbf1573ddb1f2558c1b",
		"0x5164d08bffbae85dc69385d9c575ca7d069c90d48a3f41c71ea2e4263dbd93a4",
		"0x4bb899829e70934f6e6159603961c875b9221a036cd87c4666b9623cd8e8dc9a",
		"0x5aaa12de101f9b75dd4ec9612d99993c6079773cb4cc8a5fb6c21d7b6b25b11c",
		"0x21fb36a20a703a33bc28729a392ccbd5ff2aa381afcdb1afd5bdacb811728470",
		"0x65f78557fb13e016275b55b82dad5b628edd314916be6d229c413de8b49e9c72",
		"0x29e1fdc6cf8e39ead737f611311acd9581edbf470fd52633f603cd9ca020ec33",
		"0xc1debe1a4bcb0e681cdb0a6eb0cad253a62f237aa28d4f84129eb0f5858390cd",
		"0x501fc27dbd0a0ee01febae389cc828b658f26e9a2d0777a9509dfb98c5413b27",
		"0xe2fa623b227fe49614a8d0b02d8d132f36a0f65c485c5338fce3368b32a8c142",
		"0xb94f6b71bbda55bd3c28323733b7aa0ab286e8cbd07dd1c2f3f8895776983ae7",
		"0x0776205dc04fcd110106a7d4953a60df1fa396ab0f8c29f9f383af9c5be8479b",
		"0xa22ba1151ba0422b275bc78b1c449938c420c3c6bbbecf676ac10f295cbab5cd",
		"0x6df23542b7d2a3cd1b26e34ad76764f8ff93cfcd049e25df4faa2d5cb62ffa4a",
		"0x006c22a7bc81aa8323ef4841b5e6eefa3b6265b53f00b61c5ed943f996cda021",
		"0xe28706353dda1830f7235c94e46908ad08ae10709f1b2d852321c77bdc8ca089",
		"0xec2b4fb121d080f837ec1b7619cbbf9d85b784f5dffb6979c4775975f362877a",
		"0xecfebe212b478a33b84cdbda8a106eddbac475230c6de4f82c6255a4cfe3bc07",
		"0x5d6bb875005956387738a3183a11c6fd2bf4a6cd6d1e2491fd2c6569e675fd03",
		"0x3f077a84f442a8b62320981a3d5c14892f3e387f97a179f3db7ca145dfbdfd40",
		"0x7dfb7c61db29c381c945dc2313462413280e7f4d800bb381886816197c32c7aa",
		"0x636add65634c9c452934e46b1d722b9128c2b59fcf9e4214990c70e9e5438ad1",
		"0x9204e1ffa5d9ca50986dbe3bed13994800ed11a76c51f8cf642ea880bb441ccb",
		"0x31d32cf370e30dc1fec91e616885221791d58e63ee5a15bc9c318cc7c5256402",
		"0x505efefebb226534f0046969eb57359909f8ec2145a16f6ba009371df2734b2b",
		"0x61f524c31841e8e6ff7ad3ae16f4162b568d77ec966c3fd2b3b2a466ccbb48f6",
		"0x0e7ca5c40aebc5d913f50a03c3dbd4bc24ee99e91efb9c1729588e1749475d91",
		"0x5f2c68f3f46b697ca4f22b7ea79df7522cccc96eed5b70a2d9f08bc372eb61f2",
		"0xfc6216b461b2b037408d805d86b22843b4553f1aef7e8acd0355b87eb8967b2e",
		"0xdccf535e83c067d7f1aab694b57f2c72895e89e772ce66b8702e59c7ec9a78fa",
		"0x107b093fa1c5534c9e4672b907a200d0dd1950ff855fa941e3069f669392fab6",
		"0x009344adab32cb28a110db25256d45bb177491eb0d75a4592a1830ed7e464643",
		"0x4b07c6199ce82823ac5b8463a7dae2e91680d6855df8ffe362d4703e6ccf0c82",
		"0x0b9cc54a3c6603d9c1811ac279cd5c4829c7046294341dd825e6075f4bfbaddd",
		"0xbecedc1520e0b11eaf625a0de4b1f645996dc8bf7446029ab8a82bbfc089eb59",
		"0x85328ac5ca882269b44e27030df780322e4bbe115ea9b25de404ec99f24ce572",
		"0x34848eac444746d268ecca75b352bc5e7b278f0cb6c6a48cb159de2586d04957",
		"0x4d694772918b0501cad89fa1b66b60f6a2f5111658fa9a7c11f8d77185300ddf",
		"0x3ac9372eee40ed85645a6bc506fa0254e3ffca3998c0912b6fa8a6deacf16935",
		"0x00401fa98e36e85da3f5f3df709aca6c30350b2e87dfc18b37d63576c6759c1b",
		"0x90a49e41b8a6db2cfde112069e22c7eab419733d2b81f503f1fbfec5da4053ee",
		"0xa3e0de1a5d7dbf314069c6e80d85744dc466d380af599063ec11054028ac4586",
		"0x17e95cb4f4ea5dc7cdc3829a4e284febcd9599dfc22feed62445a1f634993e26",
		"0x438a764b4baec34d5de69023f3224a83784d64f04e6e7f044f50f359609cd67c",
		"0xdb8f288df068a99893e822bfe2e8b9fbb074dae78ce4179e5746297fb66a7c3a",
		"0xac0e93cea8bee92ddfd505a56c80bb43b71bfd10698d3d7088e20a843dfff63e",
		"0x4dc4e736aed00ddf409bc539b028711fd457bc56e78de0702feb56213213ffae",
		"0xc32e463db939e8ec5b37bceb1b168c4e274d98a1536aa6caefc4d4ab7add6fda",
		"0xda48dc94e293efc7ac5ca42df4c486830dcdf7c3f4350a2f02ae6338efd15953",
		"0x9730f1be423bb3fb24bcfb29defe43278d6aa8639dea04edc4463f41c4bab559",
		"0x1d1c57411aa24a664e4a8c3b07a8c13cc63119b55651b38551a889aeba5e07cd",
		"0xd3fbedd8b4af83e7cfc7d4e6e5cb6533384735c256029d16a5fd93b4f7f48fd2",
		"0x999c33ff237eab850203ac5feda2e21d15be70d8879eaeac0093045393b4f187",
		"0x82f6fe34df24b8099f1cbbf10bf8212c2489d8da38ac326be8430449def7bb29",
		"0xdd5f201f08e613f188bde35f924a8f9956b503cbf0c9bc1352ab3b36bb09165c",
		"0xb4f5e3c56a09a0ae3ed9615314023c88e40b47a5c214a199b148959fa867d162",
		"0x16366ad8d21d0190802c80cc20f53ec35581ad16498049c1be9de9db14e9084c",
		"0xcdd5a3316db42fb053eb60a6ae96d7ad2eb59649a4c7bdc4175b1ae3ab05f99a",
		"0x5b0a82484421c3f15cb82dca47a86f443b306dd0b8c7d91be9beca02858dc742",
		"0xbbd7ddca9dd72af91f408e7e7008b43f295dae37c77adc6cfb3af4510959762c",
		"0xdd7e5f1f3bb825ce019f7d1ad5a80b2f68001ea65bd391d3eaa80d39011c745d",
		"0xd331a67194327329a492a5ae987b9673d36bbeebaa3672c8996d41e9a59f9f4b",
		"0xb99e08fdef4561960ecef9b094ac63c6026310d00637c6cb956d49467ed6a645",
		"0xaa656bdddd5135cd9fc9bc027ad4548c7b9746781ec13ed5da5569a3f66a1ef8",
		"0x01ebd7e0dad0c374e55d101494ea55f493c5c564104c662b6a593ca32fa60ad4",
		"0xd828647afcfb16622410faf46528f04c93c4d8de85d07f16ac3ceb3b3f95e195",
		"0xad030a901c15b8b176c9e304383a7b69e5bd902b12ed1c710547626586677f82",
		"0xb5d36035b78762d0ef268b18cdafe76606f19f9ba8fc64f4d8f9ac84df7e8697",
		"0x77f7989ccf1528d3049facd8269ccad1b31200546b9f18f0058a76586ad2257d",
		"0xbb952c6c79fbc7928c50270cc3c9010c71f3726badbaec1a2470693bdf23db13",
		"0xed5ad38b829d5b0098b79d338255de6a90acd60f1f72401076da89c219a063e0",
		"0xd5e3b63f6d1255225fe1de22152bcc775d10a2e211c0965a880226b48526c093",
		"0x5e8a532b7cc07847c428b62b79def81572fa2ed9c32544e4398fc3eb0437cc5e",
		"0x467e535642a28b8c56ad7cc9c77a28e95a24efc69a0f9e276c42ae382302d08e",
		"0x444299992f6c7db16ffb54a450ab1f18bc05eb93eb86cff12f166019252901bd",
		"0xbf5d1957155e271159ce7b3c866918f5afe0e5d56aa63e02e39935f8e0e668e7",
		"0xdb28b5385e96cde077e187b39df905eaf884d034910f6b8062984428f0e13838",
		"0xb5fd3c31852db20b7423e98c93ede7fb07d315b54627b382fa4a20122c236478",
		"0x7aa2b718084b542d1a0399fb02b9742447d9642aa7ec01527cd9b5d400c7f63a",
		"0x899e4286a014ff1cd9a00f95b5f02823ec64d35dde4961b5032b27a9e0918883",
		"0x361c80fa06cd2647cedf5531c1d64c305f913a6161044716fe4190516ec3081c",
		"0x27e108694a0254af878e9b6399bce3e70e9231eab687dd471aabf13345b14b19",
		"0x6a62ac86180d5b9181f6f152fad3c53582bec13e304a7de2668d5ae82819b750",
		"0x8b88c7ec8b123ccd1200a5941f08b4fd565a7899d35c170b677e0e04578e8c1e",
		"0x455ddd7ef97ec45a1c820422deb54a14e36163d2296672c1b3b7d27fbb2da488",
		"0xf9ba8c55216377e3c764c96e339c6dc64527c0eebd107060f41ec70682aab41a",
		"0x912838e937c2c9e34f95023e1a91b22295a9fd209cc3e860c9d7a3096c4b1ba2",
		"0xa492cef26f358c924b847605d95790ca9e641fefbdec6d6fa806aade8fd02ba4",
		"0x55ec949d69d952b04c900f1e5e0ba1aab99c7258b171211b4d99b2f40a071615",
		"0x621a7c46f222aad941551ef61e9aa4ca3e44783373d58727f06ef10c281a07e8",
		"0x7c262d240afbf4619971defefd7e3c36e273e594da8efb65bad9d33746ce2641",
		"0x45d34929beec87380d161289e4b0b98467e690a001c6203b35f89d8d625b61fc",
		"0x81bc17875ef16e520f65d81ffc77dd36ec73c21ba1373fe4efa1197e42e8e052",
		"0x7fb4077dc8e1c483627d807b6f4fd06f8307b5a5eaf0645d6e7bff4c7867abfe",
		"0x70952b334095f0b0f441a878118fe177bc02c65c8700080205a40ae0d2ec18ea",
		"0x24dc7f96dae6cf25fc1e338fd6cbbb0aab4a989e142449bd31137127d2631361",
		"0xa840ff38ec35d4962a05224dd8f14dc833b99de34d19e670841f763d1598b92f",
		"0x794884e89533f8997aff72cb86c9f15dfc7ccffd44295d9c85b48f79787aa73c",
		"0xd5c57ac01bdb11fd39c6e0b99423245e115a6f01f7b052d14fef563dd76bbe20",
		"0x470f9169e8a1b9dfa1f5e67352d309812e73288bb62b1fead6a41370ff9768fd",
		"0x9acd19c42a9fdd8491da549142e91f4856651750dee86323094e53062cc2d8a4",
		"0x2e4d3aa8a340527bf0a41a2cb77161bb8a47f466c3d33ae8eee98ed033b03849",
		"0x3ab8ad7c49ea6724285f68b359f9f4fc57ceb157ecbf8b192fbaafe81080d285",
		"0xa859a3ccf6f614eb4021f944c7b921a24182613ab0a439eeaefa4586e6aa0236",
		"0x2131ff15038c7a3a81777e0e8bf0733056ae7c2387851d7ea531bf3daa4ce50a",
		"0xf9eb8746cbbcf7d81c51d528aa274b5e9d27978bd128c22663b2047cbe195c73",
		"0x9c21fce1a908545e8b2f04f63c7a8b6996e27f3f34b23bbdbc07798f8aad79ff",
		"0x242f41a719816baba345d45885f9e6f1fa7237607b9a7af504129fabfcd68dd6",
		"0x1d242decb5a4b814bc177ac319325f7cc15c8a1f846a93d9b8f32f7832cfb64f",
		"0x11f953651801f6a8e17a6acd597b66d6ccc6404115f0d7a48e0271fdb3012699",
		"0x13bd51c0becfa513e8f0fbe88b57fe41156a9376c67417daeb8d71138dc397e0",
		"0x4fc7bc1c4834defb1b8ece3690433e1496e72647d70fcb2d0eb19f13b1521e72",
		"0x015c74598878b4fc11ce414c7460cdc3bf0fdd21884e05a9e13969dca251f771",
		"0x4d8ff34926cd6685dc52a85a740f07b189f2e79f5847f192b139c01afd7fbb9c",
		"0x8d194a633d850da194fdd700073847922a4622ec8a332fc7dfef15b0be23cf90",
		"0x14cdf8fa0bc39a355726d64eaad00af892d1a9937a76365af878d99c6578fa29",
		"0x727172556ff31a12bbf864b9ae8dbc970bd496762a5705ea4939dfcc0c0ccd14",
		"0x283930de85a1e85697ca4517e7055dc7dbfb43975837fc2cd16e538136e1c23b",
		"0xc9b78ee26f020d48dff7a090962b23e0c7b5559369ac9173b3d4d8fe12760ee5",
		"0x852b9d6046b4b2354de2aa87a03d1a78635c1f645d5d73a266bd3056aa6180d5",
		"0x16762314fd7c4c3af47088103743af7782e24272d1aa21abe83ef806328fcf20",
		"0xb8ba12f814705bd449fb96df9894972dd6238dd92ab55bb7ff96f76757e72fc9",
		"0xa62fb534eff5598c3431dc803c11a59546de5b52a7c5673467b122fae98187cc",
		"0xa1ad5c9f0b0a5c64f7938714e154c3debfd7438ba069eccad97bc4d511921e6e",
		"0x3ce1e7be38dea17527c6c6a40965538742a074496de999c4d6cb0d3d6e55f36b",
		"0x974a505878acaeae7e0b298a45118e5ad56f2d16853fea7cdb151f0715f29910",
		"0xfdf533005fcaa58f663c9d388d610e82c2a320a38116ab89c4a1edc5b650228f",
		"0x50da7109f228e58f2784ff4af4b0478075a7bb1d6100d58a133291cc89255471",
		"0x168a9c1ece253d5fd1c81d5c350b759c99e6eb3cbf8bcd9b913f2aac416c6bce",
		"0xb8c6ae753c2484167eca32dc57389b8b7e3b5e88c982ec41d0fa820eaf131576",
		"0x7cb81dbe6ce687b89104c5c4317844d8d966e17fc51f0851a771963d598eb7c1",
		"0xec0065cdfb97d5a3cc735b7e30e80f9325ea8e58676b8cb66666fcfd5c9d3e0d",
		"0x4cf029dd9750ccc21b26f95a588954a69866158c5979e0a364abf1fd6ceae265",
		"0x929901cfa5f0d88ebe529a258cfa05915045df796f06f8e68ed514d7af91d87b",
		"0xb0b11f437da8bc85dfb00454102b7f42190110141e564abed8342921ea63186a",
		"0x80e218d9ad8c1d6a2f546a13826043ead05c2596c4eb9399850c50c70ee1a398",
		"0x4e6dc3d3c380d8b867d0af85fe94e1754069ad18e1e69f763cf29326dcda1698",
		"0x625f0598aa5633860d516af1ab0b8993dfc86823451ebb8df3c5df3a7765953d",
		"0xfa3d7e9efa78e628b7109d0c7d730057153cd519b07978deee1e43e959ddbfeb",
		"0x854c19e049c3afb443d0b4368c453d8716311d283daff077acc41c54fa6e8d0e",
		"0x199716858f7ce1eff9e3629e2de248a503c6b724b44bfa63b64edcf19d38a119",
		"0x8410d63daae614aa4fa3c0e9f4d2f7da5641087179d35e6c96ffd66b5724f490",
		"0x267ff909220ef07a23ed3eb82025e0dff325b5a8abdafc3099842c71868d4663",
		"0x4229952d05a3924e5ceb052668aa5a1453ac224bb71fc64ea7ab3ea8e75af0d4",
		"0x4801388980d9cab4cc72f16d4cf4128aadcee43d93b60b9be0ee722fdc64c801",
		"0xde08a6f43ac5dac1a2b6f2dfcea4ae10fa40c75d895d1360fbd88478d188d2b3",
		"0x85504f37e132914499a4c350bc11b876d14de2c1715965475534651c2f815ec1",
		"0x2ea9729efaa057bedddfd8fed34c41d5da05c8f54a4d04bf5aa30ba5a3f3b05e",
		"0x54a50bb40d1369320d303191914244f4c083acbf8ee195037ecf2e2157f51390",
		"0xcb4e5de237025a0d18e1aacc818ae0eb88d01e13cd0c28315985a2c7d75ecaf5",
		"0xd5032580018d6828457c56113959a2d48a1d4b1484e6b84cdd48000df53ce838",
		"0x0683119e136ba9af7b1b9091e517462b8021f3fb912c00c331f04eae4a94c511",
		"0x5c94b9c30c6ab66e402ced1bb5bb86fb6e9391ead062fda2ca7fc409c47331b3",
		"0x231e6a27d8e93ba2930a70747da599526a001aa944bf6a130180ec2ed7a0502e",
		"0x641e5d85aa4f7feff25e87dd792b9cc2dca3be80d502c109877918319b7104e9",
		"0xa78ae4f626c87e872e6806d46a3a69f6adf30d95fd113c4e262b4ff920cb4a0c",
		"0xa39ff010749c71076e847a364649a5b90c84ed02717e1fec9d8ca90796ef4422",
		"0x18188f22b33f0ccee4c9aa33100e6e81e6cb4f66d2acdcafbef6ce38e3cfcf5a",
		"0xcf978ffe6f0f846a4b58a2443df33e26f5d3ac397811416e9b33f3f0186db300",
		"0x4ff39b891ba7c53a19df628f178390d4270bc219822b5f720cf16c64eff73345",
		"0x5d8a56ac076eca2ce269a7c74fa29dfbdda95f19a38302fc02e726b9a4f4bc16",
		"0xf5ce9d605428bcd22f77425f3d81b324aafd50692b6c84a3bb8998e3a709d431",
		"0xac48cf3541d1ca2c21b539d7959f7fe38b86d02570b0ef5d18d41ad736a43ff3",
		"0xbc643f435efb57c17d47b097eb6c2b65219c6419cae70249bf4829fc657cd8f1",
		"0x4bf8bf2b87e0b17233e959e57b89d92378545af394fa16972e98506320c0db6a",
		"0x92be41e4431520952dab6be7fb0f8daa206eb57a69fb8bc6f993b94a25af61e7",
		"0x4582198b4a38e31a37608054dbacf7454a8f7fb925f032500079fe8a1e2a9af7",
		"0xf561abb71675c899bfd977b307356314fb473abfe862b05f862dec00ee1c5812",
		"0x67b434710950cab6b199ce0ea977ff57a4a5822d758bb8bb45dc553e5129db72",
		"0xfd780478a407d506b2ae25115fb193399ab4df51adf9744e61a88741bf252071",
		"0x8f6f597d1e4d647973a28233966167aeda228e5941421b90934d967f1721f613",
		"0xedbca35d3dafa4c89c8278f9478f7f3276686dde098f14f98b58d30321b47a3f",
		"0x9cf405bd634d0eb156b5265e019ffae5ed21023389f8c0ee6738632a0171baa6",
		"0xcd6779309f4d78f57e36612ad7e354d1611ad5789cafe7519550cb9a6946831d",
		"0xd54b008e224a4350702d10e30775304be26937f9f34b31149ab695dd01cb92c0",
		"0xe211b685f17b35d26c3ac62acb177729e50c4015099c70c35d9e4457500b48fb",
		"0xd47b4b03edae4536be611f76e995ed96f1a99096f1f9a2247d9e55e1628cc451",
		"0x930d280d655e4ebaffee9fe1aa437268f7431087c4257a289467ddcea5be8c47",
		"0x98579543e0ca4359b2fce2dcb818fc599866f820e0c6a734ed1eba620edcad2f",
		"0xc614e8357315052c41da505d656da985527bd2c9d8479673cfb59387e5840717",
		"0xd3823c7c7be3155cb16457714c052f4c458e689eccd30b6214c6a1e53b2248ce",
		"0xeb252d217d184ac4b6278d5692f5abf1d70bbf03e3f3a45b3caaa39617c24a94",
		"0x53df76e2e046aa9c5a2cf365cec5d350488d2ffa7b8ffb9e4dec3ea9a9b28bf3",
		"0x619b5680eee489aeae7dd9e750ec15db4a7e76de89203911cba02e3e5c5b77aa",
		"0x05742694a970c9c244f1c1ef96f81cde46db8a2f1510a76b4e3b44013a85faee",
		"0x3a1462a17e725e52dff5433f3bc791e381dbad155fc38548cf320e93bee2971c",
		"0x403b6a691e3bd44625f0b983b7ffd6d611dd7cfc60724b269f22703be965de0d",
		"0xa6a067aa9d4798fddac30748ab4a169ef9b389f29e538fe4d6670dd4279fa8a1",
		"0x47671e171cd865569eea1cc09f6c40d8577b02a13afc1648d87f7d7cf5ca0fe3",
		"0xf7f06e6bb039634833625fc5a8c39e794ba89acba4a06dfc4ecb7e85c242af8f",
		"0x350aa2262dc0cabaee85985d99a1a1510a4429ac5f68e03a238799e52f98359e",
		"0x13f84a2cb5002d65220cec437ab89377fc445d993695cc149e9529193a325c03",
		"0x79f87db8982beeb3979e75752cccddb606cd3af2ce6af846a600935f5c72a87e",
		"0x8d4a14a6b986890a6277aef74ede8a7350e88448cebbbdf9a4cbd060785a1e7e",
		"0xa8bfaabd23b0d0cd36cc5d4f657ff07e83c234db3adf29eba04f8b8d769edd82",
		"0xfd87b4af20e95251b29c57dead269ec940987360ce283fc9d96b095c1f2add32",
		"0x8ed4627c91957184e2e704844e08b1ebcbc7cfd22ac936ef6678368c5be2ade0",
		"0x52f59851ab7591878ea4cfb3e9aad47fe67ff069669e6df088018fa355431e51",
		"0x2b86b367a0bef45695f00f22dcfc61fece33b7ac5e18fe13c85fc767fcf0caab",
		"0x83ce320fa30ad5e69497227655af987ad96121a2a467c950faa59f1d88384077",
		"0x13ce634c2849b6e66aefe727c0ce245e4e34fdb540970139e8bb355598c256c7",
		"0x2ca1f1578a26c114f6cdd4f5cd4670eea1afa235cf8b146cc11f75944aa8a240",
		"0x7c15bb884987a4b5e916234e149ea7cb47987b6c562ce9bfffa81341ed35384d",
		"0xd164a73658dbb5781a55803a1859e4ac41ed6b75dbe4dead76b40e8754ded4a2",
		"0xfa4f52cda11d0bcbdac2dd3784dbc9715da8e1334a1b63342a86aac0454e8446",
		"0x2fa7a1fbe777a856ec0a7acad06a3e4edbeacc0453501ff2e311f979cd926bc6",
		"0xf6065d4f6de608671970c2f47ecd5fa6cfd4a7089dbf5ea3e1d23bfe31903dd0",
		"0x874a6e7284c76074ce187306da1fdb5d04273957e3b47f9aa886fa35bcff7c9f",
		"0x1538134621c5b2c104ff8e31146bba374310acfd0258ade039cd7ec974e46c09",
		"0x8cf754c64a9506856ad39d75c3ed8cd19c9b8e199c6823dbff6e71aeb01ca643",
		"0x6ba4655974f7f404e95be0c48e23127f34abdb29a8ab1dc94838c8c5eac79253",
		"0x9eb07b7f4b3b5455e432aedfeb22a2fbd9e1706e75e14fb5e8673f5a030ab06a",
		"0xf0aeaeaf94e6816c5617f5eb0a2a43c8cf144f7b548a9bde39862486bf91ef91",
		"0xb9528a70f6c9c99e47230451eae69ce6a1f2332a1726165e1341f0b557f5bc59",
		"0x1b570caed4885864d97eac77887104315363516939f4821f313c976a8bb337fa",
		"0x61471eaa79fdc2d4ca62f5d4865554116048d1ee789be57bd41a2374eb08ce2e",
		"0x2efbf228a172764685cc94810db1a23f12dda7b88439a10532331054a4f2954b",
		"0xabfbe169e534bbca03bfc47c1c07eb5a07d407bc26b111d7abbdaf2b9cefc72f",
		"0x99090a29736a99e77f4d6f8657267d8974a605ec782fa2492b8db557caa68c7e",
		"0x2f34c9b455462bb3a4ffb65ea87244e714c1f86ec497ea4f19927f71640a60e7",
	}

	if n == 0 || n > uint32(len(keys)) {
		panic(errors.New("validator num is out of range"))
	}

	key, _ := crypto.ToECDSA(hexutil.MustDecode(keys[n-1]))
	return key
}
