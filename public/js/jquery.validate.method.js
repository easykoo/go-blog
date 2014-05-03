jQuery.validator.addMethod("username", function (value, element) {
    var username = /^[a-zA-Z][\w]{3,18}[a-zA-Z0-9]$/;
    return this.optional(element) || (username.test(value));
}, "5-20, begin with a letter, end with a letter or number, include letter, numbers and underscore.");
jQuery.validator.addMethod("myPassword", function (value, element) {
    return this.optional(element) || value.length >= 6 && value.length <= 20 && /\d/.test(value) && /[a-z]/i.test(value);
}, "6-20, contain at least one number and one letter.");
jQuery.validator.addMethod("fileType", function (value, element, param) {
    var fileType = value.substring(value.lastIndexOf(".") + 1).toLowerCase();
    return this.optional(element) || $.inArray(fileType, param) != -1;
}, $.validator.format("Invalid picture type."));
jQuery.validator.addMethod("fileSize", function (value, element, param) {
    return this.optional(element) || element.files[0].size <= param * 1024;
}, $.validator.format("File size is great then {0}."));
// 字符验证
jQuery.validator.addMethod("stringCheck", function (value, element) {
    return this.optional(element) || /^[\u0391-\uFFE5\w]+$/.test(value);
}, "只能包括中文字、英文字母、数字和下划线");
// 身份证号码验证
jQuery.validator.addMethod("isIdCardNo", function (value, element) {
    return this.optional(element) || isIdCardNo(value);
}, "请正确输入您的身份证号码");
// 手机号码验证
jQuery.validator.addMethod("isMobile", function (value, element) {
    var length = value.length;
    var mobile = /^(((13[0-9]{1})|(15[0-9]{1}))+\d{8})$/;
    return this.optional(element) || (length == 11 && mobile.test(value));
}, "请正确填写您的手机号码");
// 电话号码验证
jQuery.validator.addMethod("isTel", function (value, element) {
    var tel = /^\d{3,4}-?\d{7,9}$/; //电话号码格式010-12345678
    return this.optional(element) || (tel.test(value));
}, "请正确填写您的电话号码");
// 联系电话(手机/电话皆可)验证
jQuery.validator.addMethod("isPhone", function (value, element) {
    var length = value.length;
    var mobile = /^(((13[0-9]{1})|(15[0-9]{1}))+\d{8})$/;
    var tel = /^\d{3,4}-?\d{7,9}$/;
    return this.optional(element) || (tel.test(value) || mobile.test(value));
}, "请正确填写您的联系电话");
// 邮政编码验证
jQuery.validator.addMethod("isZipCode", function (value, element) {
    var tel = /^[0-9]{6}$/;
    return this.optional(element) || (tel.test(value));
}, "请正确填写您的邮政编码");
// 手机号码验证
jQuery.validator.addMethod("mobile", function (value, element) {
    var length = value.length;
    var mobile = /^(((13[0-9]{1})|(15[0-9]{1}))+\d{8})$/
    return this.optional(element) || (length == 11 && mobile.test(value));
}, "手机号码格式错误");
// 电话号码验证
jQuery.validator.addMethod("phone", function (value, element) {
    var tel = /^(0[0-9]{2,3}\-)?([2-9][0-9]{6,7})+(\-[0-9]{1,4})?$/;
    return this.optional(element) || (tel.test(value));
}, "电话号码格式错误");
// 邮政编码验证
jQuery.validator.addMethod("zipCode", function (value, element) {
    var tel = /^[0-9]{6}$/;
    return this.optional(element) || (tel.test(value));
}, "邮政编码格式错误");
// QQ号码验证
jQuery.validator.addMethod("qq", function (value, element) {
    var tel = /^[1-9]\d{4,9}$/;
    return this.optional(element) || (tel.test(value));
}, "qq号码格式错误");
// IP地址验证
jQuery.validator.addMethod("ip", function (value, element) {
    var ip = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
    return this.optional(element) || (ip.test(value) && (RegExp.$1 < 256 && RegExp.$2 < 256 && RegExp.$3 < 256 && RegExp.$4 < 256));
}, "Ip地址格式错误");
// 字母和数字的验证
jQuery.validator.addMethod("chrnum", function (value, element) {
    var chrnum = /^([a-zA-Z0-9]+)$/;
    return this.optional(element) || (chrnum.test(value));
}, "只能输入数字和字母(字符A-Z, a-z, 0-9)");
// 中文的验证
jQuery.validator.addMethod("chinese", function (value, element) {
    var chinese = /^[\u4e00-\u9fa5]+$/;
    return this.optional(element) || (chinese.test(value));
}, "只能输入中文");
// 下拉框验证
$.validator.addMethod("selectNone", function (value, element) {
    return value == "请选择";
}, "必须选择一项");
// 字节长度验证
jQuery.validator.addMethod("byteRangeLength", function (value, element, param) {
    var length = value.length;
    for (var i = 0; i < value.length; i++) {
        if (value.charCodeAt(i) > 127) {
            length++;
        }
    }
    return this.optional(element) || (length >= param[0] && length <= param[1]);
}, $.validator.format("请确保输入的值在{0}-{1}个字节之间(一个中文字算2个字节)"));