/**
 * @version     $Id: terms.js 760 2011-06-24 22:36:00Z stiandidriksen $
 * @copyright   Copyright (C) 2009 - 2010 Timble CVBA and Contributors. (http://www.timble.net)
 * @license     GNU GPLv3 <http://www.gnu.org/licenses/gpl.html>
 * @link        http://www.nooku.org
 */
(function(){
var $ = document.id;

this.Comments = new Class({
    Extends: Request,
    element : null,
    form    : null,
    
    options: {
        action      : '',
        evalScripts : false,
        
        onComplete: function() {
            if (this.response && this.response.text) {
                this.element.empty().set('html', this.response.text);
                
                this.element.getElementById('rpoohcheck').set('defaultValue', '');
                
                new Comments(this.element);
            } else {
                this.get(this.url);
            }
        }
    },
    
    initialize: function(element, options) {
        options = options || {};
        this.element = document.id(element);

        var that = this;
        this.element.getElements('a[data-action]').addEvent('click', function(e) {
            if (e.target.get('data-action')) {
                e.stop();
                that.execute(this.get('data-action'), this.get('data-id'), e);
            }
        });

        this.form = this.element.getElement('form');
        this.url = this.form.getProperty('action')+ '&tmpl=';
        
        options.url = this.url;
        this.parent(options);
        
        this.form.addEvent('submit', function(e) {
            e.stop();

            var passed      = true,
                elements    = this.form.getElements('[name=comment], [name=captcha_value], [name=email], [name=username]');

            Array.each(elements, function(element) {
               if (element.value) {
                   element.removeClass('invalid');
               }  else {
                    passed = false;
                    element.addClass('invalid'); 
               }
            });

            if (passed) { this.execute('add') };

        }.bind(this));
    },
    
    execute: function(action, data, event)
    {
        var method = '_action'+action.capitalize();
        
        if($type(this[method]) == 'function') 
        {
            this.options.action = action;
            this[method].call(this, data, event);
        }
    },

    _actionReply: function(data, event)
    {
        var target = $(event.target);
        $$('[name=path]').set('value', data);
        this.form.inject(target.getParent(), 'after');
       
    },

    _actionUnpublish: function(data, event)
    {
        if (confirm('Are you sure you want to Unpublish this comment')) {
            this.options.url = [this.options.url, 'id='+data].join('&');
            this.post({url: this.options.url, enabled:0, action:'edit', _token:this.form.getElement('[name=_token]').value}); 
        } 
    },
    _actionPublish: function(data, event)
    {
        if (confirm('Are you sure you want to publish this comment')) {
            this.options.url = [this.options.url, 'id='+data].join('&');
            this.post({url: this.options.url, enabled:1, action:'edit', _token:this.form.getElement('[name=_token]').value}); 
        } 
    },
    _actionDelete: function(data, event)
    {
        if (confirm('Are you sure you want to delete this comment, please note that all replies to this comment will also be deleted, If you wish to keep replies click unpublish.')) {
            this.options.url = [this.options.url, 'id='+data].join('&');
            this.DELETE(this.form);
        }
    },
    _actionSubscribe: function(data, event)
    {
        this.post({url: this.options.url, row:this.form.getElement('[name=row]').value, table: this.form.getElement('[name=table]').value, action:'subscribe', _token:this.form.getElement('[name=_token]').value}); 

    },
    _actionReport: function(data, event)
    {
       if (confirm('Report this comment to an administrator?')) {
            this.post({url: this.options.url, comment:data, action:'report', _token:this.form.getElement('[name=_token]').value}); 
        }
    },
    _actionSpam: function(data, event)
    {
        if (confirm('Are you sure you want to mark this comment as spam?')) {
            this.post({url: this.options.url, comment:data, action:'spam', _token:this.form.getElement('[name=_token]').value}); 
        }
    },
    _actionAdd: function(data)
    {
        this.post(this.form);
    }
});
})();

window.addEvent('domready', function() {
    new Comments('lfcomments');
});