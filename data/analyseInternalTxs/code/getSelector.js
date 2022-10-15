const ethers = require('ethers')

//函数功能: 计算函数选择器
//input: funcHead 是函数头部申明，去掉""funcion"关键字，去掉参数名字，只
//          保留函数名+参数类型，去掉所有的空格
//return: 返回4个字节的字符串
function getFunctionSeletor(funcHead) {
    return ethers.utils.id(funcHead).slice(0,10)
}

//函数功能: 计算事件选择器
//input: eventHead 是事件头部申明，去掉""event"关键字，去掉参数名字，只
//          保留事件名+参数类型，去掉所有的空格
//return: 返回32个字节的字符串
function getEventSelector(eventHead) {
    return ethers.utils.id(eventHead)
}

function doMain() {
    //1)计算函数选择器
    //原函数：function allowance(address owner, address spender)
    //精简化: allowance(address,address)
    strfunc = 'ownerOf(uint256 _tokenId) external view returns (address)'
    funcSeletor = getFunctionSeletor(strfunc)
    res = strfunc + '='
    console.log(res,funcSeletor)

    //2) 计算事件选择器
    //原事件: event Approval(address owner, address spender, uint256 value)
    //精简化: Approval(address,address,uint256)
    strEvent = 'Transfer(address indexed _from, address indexed _to, uint256 indexed _tokenId)'
    eventSelector = getEventSelector(strEvent)
    res2 = strEvent + '='
    console.log(res2,eventSelector)
}

doMain()
