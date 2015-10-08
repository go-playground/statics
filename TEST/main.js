'use strict';

(function ($) {

    $.fn.JsonTablecolumn = function (options) {
        var opt = $.extend({}, $.fn.JsonTablecolumn.defaults, options);
        return opt;
    };


    $.fn.JsonTablecolumn.defaults = {
        heading: 'Heading',
        data: 'json_field',
        type: 'string',
        sortable: true,
        starthidden: false
    };

    $.fn.JsonTable = function (options) {

        var thisSelector = this.selector;

        var opts = $.extend({}, $.fn.JsonTable.defaults, options);

        var $JsonTableContainer = $('<div>', { 'data-so': 'A', 'data-ps': opts.pageSize }).addClass('json-table');

        var $visibleColumnsCBList = $('<div>').addClass('legend col-sm-3');

        var $table = $('<table>').addClass(opts.cssClass);

        var $thead = $('<thead>');
        var $theadRow = $('<tr>');

        $.each(opts.columns, function (index, item) {
            var $th = $('<th>').attr('data-i', index);
            var $inputGroup = $('<div>', {'class': 'input-group'}).appendTo($visibleColumnsCBList);
            var $spanInputGroup = $('<span>', {'class': 'input-group-addon'}).appendTo($inputGroup);
            var $cb = $('<input>', { 'type': 'checkbox', 'id': 'cb' + thisSelector + index, value: index, checked: !item.starthidden, 'data-i': index }).bind('change', opts.onHiddenCBChange).appendTo($spanInputGroup);

            var $cblabel = $('<label />', { 'for': 'cb' + thisSelector + index, text: item.heading }).appendTo($inputGroup);

            if (item.starthidden)
                {
                    $th.hide();
                }

                if (item.sortable) {
                    $('<a>', { 'class': 'glyphicon glyphicon-sort', 'href': '#', 'data-i': index, 'data-t': item.type }).text(item.heading).bind('click', opts.onSortClick).appendTo($th);
                } else {
                    $('<span>').text(item.heading).appendTo($th);
                }

                $th.appendTo($theadRow);
            });

            $theadRow.appendTo($thead);
            $thead.appendTo($table);

            var pagingNeeded = false;
            $.each(opts.data, function (index, item) {
                var $tr = $('<tr>').attr('data-i', index);

                if (opts.pageSize <= index) {
                    $tr.hide();
                    pagingNeeded = true;
                }

                $.each(opts.columns, function (cIndex, cItem) {
                    var $td = '';

                    // FORMAT INFORMATION @anthony
                    if (cItem.format){
                        item[cItem.data] = opts.cellFormat(cItem.format, item[cItem.data]);
                    }

                    if( cItem.class !== undefined){
                         $td = $('<td>').text(item[cItem.data]).attr('data-i', cIndex).addClass(cItem.class);
                    }else{
                         $td = $('<td>').text(item[cItem.data]).attr('data-i', cIndex);
                    }

                    if (cItem.starthidden) {
                        $td.hide();
                    }

                    $td.appendTo($tr);
                });
                $tr.appendTo($table);
            });

            $JsonTableContainer.append($visibleColumnsCBList);
            $JsonTableContainer.append($table);


            if (pagingNeeded) {
                var $pager = $('<div>').addClass('paging');
                for (var i = 0; i < Math.ceil(opts.data.length / opts.pageSize) ; i++) {
                    $('<a>', { 'text': 'Page ' + (i + 1), 'href': '#', 'data-i': (i + 1), 'class': 'p-link' }).bind('click', opts.onPageClick).appendTo($pager);
                }
                $JsonTableContainer.append($pager).addClass('paged');
            }
            return this.append($JsonTableContainer);
        };


        $.fn.JsonTable.defaults = {
            cssClass: 'table',
            columns: [],
            data: [],
            pageSize: 10,

            onPageClick: function () {
                var $thisGrid = $(this).parents('.json-table');

                var pageSize = $thisGrid.attr('data-ps');
                var page = $(this).attr('data-i');

                $('tbody tr', $thisGrid).each(function (trIndex, trItem) {
                    $(this).hide();

                    var pageStart = ((page - 1) * pageSize) + 1;
                    var pageEnd = page * pageSize;

                    if ((trIndex + 1) >= pageStart && (trIndex + 1) <= pageEnd) {
                        $(this).show();
                    }
                });

                return false;
            },

            cellFormat: function (obj, value){

                var decimal = 0;
                var symbolPercentage = '';

                // using Number.prototype.toFixed() from MDN. (obj.tofixed)
                /*
                 Digits

                 Optional. The number of digits to appear after the decimal point; this may be a value between 0 and 20,
                 inclusive, and implementations may optionally support a larger range of values. If this argument is
                 omitted, it is treated as 0.

                 A string representation of numObj that does not use exponential notation and has exactly digits digits
                 after the decimal place. The number is rounded if necessary, and the fractional part is padded with
                 zeros if necessary so that it has the specified length. If numObj is greater than 1e+21, this method
                 simply calls Number.prototype.toString() and returns a string in exponential notation.
                 */


                if (obj.decimal !== undefined){
                  decimal =  obj.decimal;
                }

                // Percentage : Decimal to Percentage
                if (obj.percentage === true){
                    value = value * 100;
                    symbolPercentage = '%';
                }

                return value.toFixed(decimal) + symbolPercentage;

            },

            onHiddenCBChange: function () {
                var $thisGrid = $(this).parents('.json-table');
                var columIndex = $(this).attr('data-i');

                if ($(this).is(':checked')) {
                    $('td[data-i=' + columIndex + ']', $thisGrid).show();
                    $('th[data-i=' + columIndex + ']', $thisGrid).show();
                } else {
                    $('td[data-i=' + columIndex + ']', $thisGrid).hide();
                    $('th[data-i=' + columIndex + ']', $thisGrid).hide();
                }
            },

            onSortClick: function () {
                var $thisGrid = $(this).parents('.json-table');
                var direction = $thisGrid.attr('data-so');

                $('.glyphicon-sort', $thisGrid).removeClass('s-A s-D');
                $(this).addClass('s-' + direction);

                var type = $(this).attr('data-t');
                var index = $(this).attr('data-i');

                var array = [];

                $('tbody tr', $thisGrid).each(function (trIndex, trItem) {
                    var item = $('td', trItem).eq(index);

                    var trId = item.parent().attr('data-i');

                    var value = null;
                    switch (type) {
                        case 'string':
                            value = item.text();
                            break;
                            case 'int':
                                value = parseInt(item.text());
                                break;

                                case 'float':
                                    value = parseFloat(item.text());
                                    break;

                                    case 'datetime':
                                        value = new Date(item.text());
                                        break;

                                        default:
                                            value = item.text();
                                            break;
                                        }

                                        array.push({ trId: trId, val: value });
                                    });

                                    if (direction === 'A') {
                                        array.sort(function (a, b) {
                                            if (a.val > b.val) { return 1; }
                                                if (a.val < b.val) { return -1; }
                                                    return 0;
                                                });
                                                $thisGrid.attr('data-so', 'D');
                                            } else {

                                                array.sort(function (a, b) {
                                                    if (a.val < b.val) { return 1; }
                                                        if (a.val > b.val) { return -1;}
                                                            return 0;
                                                        });

                                                        $thisGrid.attr('data-so', 'A');
                                                    }

                                                    for (var i = 0; i < array.length; i++) {
                                                        var td = $('tr[data-i=' + array[i].trId + ']', $thisGrid);

                                                        td.detach();

                                                        $('tbody', $thisGrid).append(td);
                                                    }

                                                    if ($thisGrid.hasClass('paged')) {
                                                        $('.p-link', $thisGrid).eq(0).click();
                                                    }

                                                    return false;
                                                }




                                            };

    }(jQuery));


$(function(){
    $('table').each(function() {
        if($(this).find('thead').length > 0 && $(this).find('th').length > 0) {
            // Clone <thead>
            var $w	   = $(window),
                $t	   = $(this),
                $thead = $t.find('thead').clone(),
                $col   = $t.find('thead, tbody').clone();

            // Add class, remove margins, reset width and wrap table
            $t
                .addClass('sticky-enabled')
                .css({
                    margin: 0,
                    width: '100%'
                }).wrap('<div class="sticky-wrap" />');

            if($t.hasClass('overflow-y')){
                $t.removeClass('overflow-y').parent().addClass('overflow-y');
            }

            // Create new sticky table head (basic)
            $t.after('<table class="sticky-thead" />');

            // If <tbody> contains <th>, then we create sticky column and intersect (advanced)
            if($t.find('tbody th').length > 0) {
                $t.after('<table class="sticky-col" /><table class="sticky-intersect" />');
            }

            // Create shorthand for things
            var $stickyHead  = $(this).siblings('.sticky-thead'),
                $stickyCol   = $(this).siblings('.sticky-col'),
                $stickyInsct = $(this).siblings('.sticky-intersect'),
                $stickyWrap  = $(this).parent('.sticky-wrap');

            $stickyHead.append($thead);

            $stickyCol
                .append($col)
                .find('thead th:gt(0)').remove()
                .end()
                .find('tbody td').remove();

            $stickyInsct.html('<thead><tr><th>'+$t.find('thead th:first-child').html()+'</th></tr></thead>');

            // Set widths
            var setWidths = function () {
                    $t
                        .find('thead th').each(function (i) {
                            $stickyHead.find('th').eq(i).width($(this).width());
                        })
                        .end()
                        .find('tr').each(function (i) {
                            $stickyCol.find('tr').eq(i).height($(this).height());
                        });

                    // Set width of sticky table head
                    $stickyHead.width($t.width());

                    // Set width of sticky table col
                    $stickyCol.find('th').add($stickyInsct.find('th')).width($t.find('thead th').width());
                },
                repositionStickyHead = function () {
                    // Return value of calculated allowance
                    var allowance = calcAllowance();

                    // Check if wrapper parent is overflowing along the y-axis
                    if($t.height() > $stickyWrap.height()) {
                        // If it is overflowing (advanced layout)
                        // Position sticky header based on wrapper scrollTop()
                        if($stickyWrap.scrollTop() > 0) {
                            // When top of wrapping parent is out of view
                            $stickyHead.add($stickyInsct).css({
                                opacity: 1,
                                top: $stickyWrap.scrollTop()
                            });
                        } else {
                            // When top of wrapping parent is in view
                            $stickyHead.add($stickyInsct).css({
                                opacity: 0,
                                top: 0
                            });
                        }
                    } else {
                        // If it is not overflowing (basic layout)
                        // Position sticky header based on viewport scrollTop
                        if($w.scrollTop() > $t.offset().top && $w.scrollTop() < $t.offset().top + $t.outerHeight() - allowance) {
                            // When top of viewport is in the table itself
                            $stickyHead.add($stickyInsct).css({
                                opacity: 1,
                                top: $w.scrollTop() - $t.offset().top
                            });
                        } else {
                            // When top of viewport is above or below table
                            $stickyHead.add($stickyInsct).css({
                                opacity: 0,
                                top: 0
                            });
                        }
                    }
                },
                repositionStickyCol = function () {
                    if($stickyWrap.scrollLeft() > 0) {
                        // When left of wrapping parent is out of view
                        $stickyCol.add($stickyInsct).css({
                            opacity: 1,
                            left: $stickyWrap.scrollLeft()
                        });
                    } else {
                        // When left of wrapping parent is in view
                        $stickyCol
                            .css({ opacity: 0 })
                            .add($stickyInsct).css({ left: 0 });
                    }
                },
                calcAllowance = function () {
                    var a = 0;
                    // Calculate allowance
                    $t.find('tbody tr:lt(3)').each(function () {
                        a += $(this).height();
                    });

                    // Set fail safe limit (last three row might be too tall)
                    // Set arbitrary limit at 0.25 of viewport height, or you can use an arbitrary pixel value
                    if(a > $w.height()*0.25) {
                        a = $w.height()*0.25;
                    }

                    // Add the height of sticky header
                    a += $stickyHead.height();
                    return a;
                };

            setWidths();

            $t.parent('.sticky-wrap').scroll($.throttle(250, function() {
                repositionStickyHead();
                repositionStickyCol();
            }));

            $w
                .load(setWidths)
                .resize($.debounce(250, function () {
                    setWidths();
                    repositionStickyHead();
                    repositionStickyCol();
                }))
                .scroll($.throttle(250, repositionStickyHead));
        }
    });
});