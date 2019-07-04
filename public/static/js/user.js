
layui.use(['table', 'laydate', 'form'], function(){
    var table = layui.table
        ,form = layui.form
        ,util=layui.util
        ,laydate = layui.laydate;

    laydate.render({
        elem: '#start_time',
        type: 'datetime'
    });

    laydate.render({
        elem: '#end_time',
        type: 'datetime'
    });

    table.render({
        elem: '#userList'
        ,url:'/admin/user/user_list'
        ,id: "userListRender"
        ,toolbar: '#toolbarTpl'
        ,cellMinWidth: 80
        ,cols: [[
            {type: 'checkbox'}
            ,{field:'id', title:'ID', width:60, unresize: true, sort: true}
            ,{field:'bind_id', title:'绑定ID', width:60, unresize: true}
            ,{field:'user_login', title:'用户名', width:140}
            ,{field:'create_time', title:'创建时间', width:180, sort: true, templet: function(d) {return util.toDateString(d.create_time*1000); }}
            ,{field:'last_login_time', title:'最后时间', width:180, sort: true, templet: function(d) {return util.toDateString(d.update_time*1000); }}
            ,{field:'last_login_ip', title: '最后IP', minWidth:120}
            ,{field:'score', title: '等级', width:80, sort: true}
            ,{field:'balance', title: '余额', width:80, sort: true}
            ,{field:'power', title: '权限', width:80, sort: true,templet: '#powerTpl'}
            ,{field:'user_status', title:'操作', width:200, toolbar: '#actionTpl', unresize: true}
        ]]
        ,page: true
        // ,page:{ //支持传入 laypage 组件的所有参数（某些参数除外，如：jump/elem） - 详见文档
        //         layout: ['count', 'prev', 'page', 'next', 'skip', 'limit'] //自定义分页布局
        //             //,curr: 5 //设定初始在第 5 页
        //             ,groups: 5 //只显示 5 个连续页码
        //             ,first: true //显示首页
        //             ,last: true //显示尾页
        // }
    });

    form.on('submit(sreach)', function (data) {
        var loading = layer.load(1, {shade: [0.1, '#FF0000']});

        table.reload('userListRender', {
            page: {
                curr: 1 //重新从第 1 页开始
            }
            ,where: data.field
        });
        layer.close(loading);
        return false
    })

    table.on('toolbar(userListAction)', function(obj){
        var checkStatus = table.checkStatus(obj.config.id);
        var data = checkStatus.data;
        switch(obj.event){
            case 'addUser':
                x_admin_show('添加角色','/admin/user/user_edit');
                break;
            case 'startUsers':
                changeUserStatusWithBatch(data,1);
                break;
            case 'stopUsers':
                changeUserStatusWithBatch(data,0);
                break;
            case 'deleteUsers':
                deleteUserWithBatch(data)
                break;
        };
    });

    //监听状态操作
    form.on('switch(statusSwitch)', function(obj){
        changeUserStatus(this.value, obj)
    });

    table.on('tool(userListAction)', function(obj){
        switch(obj.event){
            case 'delete':
                deleteUser(obj);
                break
            case 'edit':
                editUser(obj);
                break;
        }
    });


    function deleteUser(obj) {
        var data = obj.data;
        layer.confirm('确定删除吗', function(index){
            $.ajax({
                url: '/admin/user/user_delete',
                data: {"id": data.id,"name":data.user_login},
                type: "post",
                dataType: "json",
                success: function (ret) {
                    var message = ret.msg + ret.code;
                    if (ret.code === 0) {
                        message = ret.msg
                        obj.del();
                        layer.close(index);
                    }
                    layer.msg(message, {icon: 1, time: 1000}, function () {
                    });
                }
            });
        });
    }
    function deleteUserWithBatch(data) {
        var users = [];
        var users_name = [];
        data.forEach(function(value,i) {
            users.unshift(value.id)
            users_name.unshift(value.user_login)
        });
        if (users.length >0 ){
            layer.confirm('确定删除吗', function(index){
                $.ajax({
                    url: '/admin/user/user_delete',
                    data: {"id": users,"name":users_name},
                    type: "post",
                    dataType: "json",
                    traditional: true,
                    success: function (ret) {
                        var message = ret.msg + ret.code;
                        if (ret.code === 0) {
                            message = ret.msg
                            table.reload('userListRender')
                        }
                        layer.msg(message, {icon: 1, time: 1000}, function () {
                        });
                    }
                });
            });
        }else{
            layer.msg("未选中目标", {icon: 1, time: 1000}, function () {});
        }
    }

    function editUser(obj){
        var data = obj.data;
        x_admin_show('添加角色','/admin/user/user_edit?user_name='+data.user_login+'&user_id='+data.id);
    }

    function changeUserStatus(id,obj) {
        var status = 0
        if (obj.elem.checked) {
            status = 1
        }
        $.ajax({
            url: '/admin/user/user_status_change',
            data: {"id": id,"status":status},
            type: "post",
            dataType: "json",
            success: function (ret) {
                var message = ret.msg;
                if (ret.code < 0) {
                    message += ret.code;
                    obj.elem.checked = !obj.elem.checked;
                    form.render();
                }
                layer.msg(message, {icon: 1, time: 1000}, function () {});
            }
        });
    }

    function changeUserStatusWithBatch(data, status) {
        var users = [];
        data.forEach(function(value,i) {
            users.push(value.id)
        });
        if (users.length >0 ){
            $.ajax({
                url: '/admin/user/user_status_change',
                data: {"id": users,"status":status},
                type: "post",
                dataType: "json",
                traditional: true,
                success: function (ret) {
                    var message = ret.msg + ret.code;
                    if (ret.code ===0) {
                        message = ret.msg
                        table.reload('userListRender')
                    }
                    layer.msg(message, {icon: 1, time: 1000}, function () {});
                }
            });
        }else{
            layer.msg("未选中目标", {icon: 1, time: 1000}, function () {});
        }
    }
});
