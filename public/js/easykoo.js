String.prototype.endWith = function (s) {
    if (s == null || s == "" || this.length == 0 || s.length > this.length)
        return false;
    if (this.substring(this.length - s.length) == s)
        return true;
    else
        return false;
    return true;
};

var formatTime = function (timeString) {
    var date = timeString.substr(0, 10)
    var time = timeString.substr(11, 8)
    return date + " " + time;
};

var changeLanguage = function (lang) {
    $.ajax('/language/change/' + lang, {
        dataType: 'json',
        type: "GET",
        success: function (data) {
            if (data.success) {
                location.reload();
            }
        }
    });
}

var filterSqlStr = function (value) {
    var sqlStr = sql_str().split(',');
    var flag = false;

    for (var i = 0; i < sqlStr.length; i++) {
        if (value.toLowerCase().indexOf(sqlStr[i]) != -1) {
            flag = true;
            break;
        }
    }
    return flag;
}

var sql_str = function () {
    var str = "and,delete,or,exec,insert,select,union,update,count,*,',join,>,<";
    return str;
}

var cutoff = function (content) {
    var cutoffLine = "----------"
    var index = content.indexOf(cutoffLine);
    if (index > 0) {
        content = content.replace(cutoffLine, "");
        var pre = content.substr(0, index);
        var preIndex = pre.lastIndexOf('</p>');
        if (preIndex > 0) {
            preIndex += 4;
        } else {
            preIndex = pre.lastIndexOf('</div>');
            if (preIndex > 0) {
                preIndex += 4;
            } else {
                return content;
            }
        }
        var pre = content.substr(0, preIndex);
        var nex = content.substr(preIndex, content.length);
        return pre + cutoffLine + nex;
    }
    return content;
}

var direct = function () {
    var winHeight = $(window).height()

    var $top = $('#goTop');
    var $bottom = $('#goBottom');
    var side = $('#side').offset().left;
    var width = $('#side').width();
    var pos = side + width + 25;
    $top.css({
        "left": pos + "px",
        "top": winHeight / 2 - 23 + "px",
        "width": "45px",
        "height": "45px",
        "position": "fixed",
        "opacity": .4
    })
    $bottom.css({
        "left": pos + "px",
        "top": winHeight / 2 + 23 + "px",
        "width": "45px",
        "height": "45px",
        "position": "fixed",
        "opacity": .4
    })
    $(window).scroll(function () {
        var scroll = $(this).scrollTop()
        if (scroll > 0) {
            $top.removeClass("hidden");
        } else {
            $top.addClass('hidden');
        }

        if (scroll + winHeight == $(document).height()) {
            $bottom.addClass('hidden');
        } else {
            $bottom.removeClass("hidden");
        }
    });
    $top.on("click", function () {
        $('html, body').animate({scrollTop: 0}, 300);
        return false;
    })
    $bottom.click(function () {
        $('html, body').animate({scrollTop: $(document).height()}, 300);
        return false;
    });
}

var setupPage = function (section, pageNo, totalPage) {
    var html = "";
    if (totalPage > 1) {
        html += '<ul class="pagination">';
        if (pageNo > 1) {
            html += '<li><a class="page_prev" href="javascript:goToPage(' + 1 + ')"><i class="fa fa-angle-double-left"></i></a></li>'
            html += '<li><a class="page_prev" href="javascript:goToPage(' + (pageNo - 1 < 1 ? 1 : pageNo - 1) + ')"><i class="fa fa-angle-left"></i></a></li>'
        } else {
            html += '<li class="disabled"><span class="fa fa-angle-double-left"></span></li>'
            html += '<li class="disabled"><span class="fa fa-angle-left"></span></li>'
        }
        for (var i = 1; i <= totalPage; i++) {
            if (i == pageNo) {
                html += '<li class="active"><span>' + i + '</span></li>';
            } else if (i < pageNo - 5) {
                if (i == 1) {
                    html += '<li><a href="javascript:goToPage(' + i + ')">' + i + '</a></li>';
                } else if (i == pageNo - 6) {
                    html += '<li><span>...</span></li>';
                }
            } else if (i > pageNo + 5) {
                if (i == totalPage) {
                    html += '<li><a href="javascript:goToPage(' + i + ')">' + i + '</a></li>';
                } else if (i == pageNo + 6) {
                    html += '<li><span>...</span></li>';
                }
            } else {
                html += '<li><a href="javascript:goToPage(' + i + ')">' + i + '</a></li>';
            }
        }
        if (pageNo < totalPage) {
            html += '<li><a class="page_next" href="javascript:goToPage(' + (pageNo + 1 > totalPage ? totalPage : pageNo + 1) + ')"><i class="fa fa-angle-right"></i></a></li>';
            html += '<li><a class="page_prev" href="javascript:goToPage(' + totalPage + ')"><i class="fa fa-angle-double-right"></i></a></li>'
        } else {
            html += '<li class="disabled"><span class="fa fa-angle-right"></span></li>'
            html += '<li class="disabled"><span class="fa fa-angle-double-right"></span></li>'
        }
        html += '</ul>';
        section.html(html);
    }
}

var initTag = function () {
    $(".tag").hover(function () {
        $(this).find(".badge").removeClass("hidden");
        $(this).parent().css("padding-right", "8px");
    }, function () {
        $(this).find(".badge").addClass("hidden");
        $(this).parent().css("padding-right", "15px");
    })
}