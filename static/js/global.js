$(function() {
    $(document).foundation();
		
	// Hide any messages after a few seconds
    hideFlash();
});

function hideFlash(rnum)
{    
    if (!rnum) rnum = '0';
    
    _.delay(function() {
        $('.alert-box-fixed' + rnum).fadeOut(300, function() {
            $(this).css({"visibility":"hidden",display:'block'}).slideUp();
            
            var that = this;
            
            _.delay(function() { that.remove(); }, 400);
        });
    }, 4000);
}

function showFlash(obj)
{
	if (!$.isArray(obj)) {
		return;
	}
	
    $('#flash-container').html();
	
    $(obj).each(function(i, v) {
        var rnum = _.random(0, 100000);
		var message = '<div data-alert id="flash-message" class="alert-box-fixed'
		+ rnum + ' alert-box-fixed alert-box '+v.Class+'">'
		+ v.Message + '<a href="#" class="close">&times;</a></div>';
        $('#flash-container').prepend(message);
        hideFlash(rnum);
    });
}

function quickFlashSuccess(message) {
	var flash = [{Class: "alert-box success", Message: message}];	
	showFlash(flash);
}

function quickFlashAlert(message) {
	var flash = [{Class: "alert-box alert", Message: message}];
	showFlash(flash);
}

function quickFlashWarning(message) {
	var flash = [{Class: "alert-box warning", Message: message}];
	showFlash(flash);
}

function quickFlash(message) {
	var flash = [{Class: "alert-box", Message: message}];
	showFlash(flash);
}

function approvePhoto(userid, picid) {
	$.get("/admin/approve/"+userid+"/"+picid, function(data) {
		showFlash(data);
	});
}

function rejectPhoto(userid, picid) {
	// Get the rejection note
	var note = $("#note_"+userid+picid).val();
	
	// Show flash if no note exists
	if (note.length < 1) {
		quickFlashAlert("You must type in a rejection note.");
		return
	}
	
	// Empty the text
	$("#note_"+userid+picid).val("");

	$.get("/admin/reject/"+userid+"/"+picid, {note: note}, function(data) {
		showFlash(data);
	});
}

function unverifyPhoto(userid, picid) {
	$.get("/admin/unverify/"+userid+"/"+picid, function(data) {
		showFlash(data);
	});
}

function deletePhoto(picid) {
	$.get("/profile/photo/delete/"+picid, function(data) {
		$("#photo"+picid).remove();
		showFlash(data);
		if ($("#photo-container").children().length == 0) {
			var baseuri = $("#BaseURI").val();
			window.location.href = baseuri;
		}
	});
}

function downloadPhoto(picid) {
	$.get("/profile/photo/download/"+picid, function(data) {
		showFlash(data);
	});
}

function verifyUsername() {
	var username = $("#username").val();
	var site = $("#site").val();
	var baseuri = $("#BaseURI").val();
	
	if (site.trim().length < 1) {
		quickFlashAlert("You must select a dating website.");
		return
	}
	
	if (username.trim().length < 1) {
		quickFlashAlert("You must type in a username.");
		return
	}	
	
	window.location.href = baseuri + "public/" + encodeURIComponent(site) + "/" + encodeURIComponent(username);
}

$("#username").keypress(function(e) {
	if (e.which == 13)	{
		verifyUsername();
		return false;
	}
});